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
	Handler(ctx context.Context, msg *kafka.Message) (interface{}, error)
	Batch(ctx context.Context, ch chan interface{}) error
	Name() string
}

type HandlerFunc func(ctx context.Context, msg *kafka.Message) (interface{}, error)
type BatchFunc func(ctx context.Context, ch chan interface{}) error

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
	_, err := e.maker.Make(h.Name())
	if err != nil {
		return err
	}

	e.handlers[h.Name()] = h
	return nil
}

func (e *Process) ProvideRunGroup(group *run.Group) {
	if len(e.handlers) == 0 {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())

	for _, one := range e.handlers {
		msgChan := make(chan *kafka.Message, 200)
		batchChan := make(chan interface{}, 200)

		group.Add(func() error {
			reader, err := e.maker.Make(one.Name())
			if err != nil {
				return err
			}
			for {
				select {
				case <-ctx.Done():
					return nil
				default:
					message, err := reader.ReadMessage(ctx)
					if err != nil {
						return err
					}
					if len(message.Value) > 0 {
						msgChan <- &message
					}
				}
			}
		}, func(err error) {
			cancel()
			close(msgChan)
			return
		})
		group.Add(func() error {
			for {
				select {
				case msg := <-msgChan:
					v, err := one.Handler(ctx, msg)
					if err != nil {
						return err
					}
					if v != nil {
						batchChan <- v
					}
				case <-ctx.Done():
					return nil
				}
			}
		}, func(err error) {
			cancel()
			close(batchChan)
			return
		})

		group.Add(func() error {
			err := one.Batch(ctx, batchChan)
			if err != nil {
				return err
			}
			return nil
		}, func(err error) {
			cancel()
			return
		})

	}
}
