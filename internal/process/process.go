package process

import (
	"context"

	"github.com/DoNewsCode/core/di"
	"github.com/DoNewsCode/core/otkafka"
	"github.com/go-kit/kit/log"
	"github.com/oklog/run"
	"github.com/segmentio/kafka-go"
)

type Process struct {
	maker    otkafka.ReaderMaker
	handlers map[string]H
	logger   log.Logger
}

type H interface {
	// Handler decode the *kafka.Message and filter
	Handler(ctx context.Context, msg *kafka.Message) (interface{}, error)
	// Batch processing of decoded or filtered data
	//
	// Note: if sequence is necessary, make sure that batch workers count is one
	Batch(ctx context.Context, data []interface{}) error
	// Info define the topic name and some config
	Info() *Info
}

type HandlerFunc func(ctx context.Context, msg *kafka.Message) (interface{}, error)
type BatchFunc func(ctx context.Context, data []interface{}) error

type in struct {
	di.In

	Hs     []H `group:"H"`
	Maker  otkafka.ReaderMaker
	Logger log.Logger
}

type Out struct {
	di.Out

	Hs []H `group:"H,flatten"`
}

func NewProcess(i in) (*Process, error) {
	e := &Process{
		maker:    i.Maker,
		logger:   i.Logger,
		handlers: map[string]H{},
	}
	for _, hh := range i.Hs {
		if err := e.addHandler(hh); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *Process) addHandler(h H) error {
	name := h.Info().name()
	_, err := e.maker.Make(name)
	if err != nil {
		return err
	}

	e.handlers[name] = h
	return nil
}

type batchInfo struct {
	message *kafka.Message
	data    interface{}
}

func (e *Process) ProvideRunGroup(group *run.Group) {
	if len(e.handlers) == 0 {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())

	var g run.Group
	msgChs := make([]chan *kafka.Message, 0)
	batchChs := make([]chan *batchInfo, 0)

	for name, ooo := range e.handlers {
		one := ooo

		msgCh := make(chan *kafka.Message, one.Info().chanSize())
		batchCh := make(chan *batchInfo, one.Info().chanSize())
		msgChs = append(msgChs, msgCh)
		batchChs = append(batchChs, batchCh)

		reader, _ := e.maker.Make(name)

		for i := 0; i < one.Info().readWorker(); i++ {
			g.Add(func() error {
				for {
					select {
					case <-ctx.Done():
						return nil
					default:
						message, err := reader.FetchMessage(ctx)
						if err != nil {
							return err
						}
						if len(message.Value) > 0 {
							msgCh <- &message
						}
					}
				}
			}, func(err error) {

			})
		}

		for i := 0; i < one.Info().handleWorker(); i++ {
			g.Add(func() error {
				for {
					select {
					case msg := <-msgCh:
						v, err := one.Handler(ctx, msg)
						if err != nil {
							return err
						}
						batchCh <- &batchInfo{message: msg, data: v}
					case <-ctx.Done():
						return nil
					}
				}
			}, func(err error) {

			})
		}

		for i := 0; i < one.Info().batchWorker(); i++ {
			g.Add(func() error {
				err := e.batch(ctx, reader, batchCh, one.Batch, one.Info().batchSize())
				if err != nil {
					return err
				}
				return nil
			}, func(err error) {

			})
		}
	}

	group.Add(func() error {
		if err := g.Run(); err != nil {
			return err
		}
		return nil
	}, func(err error) {
		cancel()
		for _, ch := range msgChs {
			close(ch)
		}
		for _, ch := range batchChs {
			close(ch)
		}
	})

}

func (e *Process) batch(ctx context.Context, reader *kafka.Reader, ch chan *batchInfo, batchFunc BatchFunc, batchSize int) error {
	var data = make([]interface{}, 0)
	var msgs = make([]kafka.Message, 0)
	for {
		select {
		case v := <-ch:
			data = append(data, v.data)
			msgs = append(msgs, *v.message)
			if len(data) >= batchSize {
				// do something, such as batch insert to db
				if err := batchFunc(ctx, data); err != nil {
					return err
				}
				if err := reader.CommitMessages(context.Background(), msgs...); err != nil {
					return err
				}
				data = data[0:0]
				msgs = msgs[0:0]
			}
		case <-ctx.Done():
			for v := range ch {
				data = append(data, v)
				msgs = append(msgs, *v.message)
			}
			if len(data) > 0 {
				// do something
				if err := batchFunc(ctx, data); err != nil {
					return err
				}
				if err := reader.CommitMessages(context.Background(), msgs...); err != nil {
					return err
				}
				data = data[0:0]
				msgs = msgs[0:0]
			}
			return nil
		}
	}
}
