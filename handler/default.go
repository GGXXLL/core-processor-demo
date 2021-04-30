package handler

import (
	"context"
	"encoding/json"
	"github.com/GGXXLL/core-kafka/entity"
	"github.com/GGXXLL/core-kafka/internal/process"
	"github.com/go-kit/kit/log"
	"github.com/segmentio/kafka-go"
)

type defaultHandler struct {
	logger log.Logger
}

func newDefaultHandler(logger log.Logger) process.Out {
	return process.Out{Hs: []process.H{&defaultHandler{logger}}}
}

func (h *defaultHandler) Name() string {
	return "default"
}

func (h *defaultHandler) Handler(ctx context.Context, msg *kafka.Message) (interface{}, error) {
	e := entity.Example{}
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (h *defaultHandler) Batch(ctx context.Context, ch chan interface{}) error {
	var l = make([]*entity.Example, 0)
	for {
		select {
		case v := <-ch:
			if e, ok := v.(*entity.Example); ok {
				l = append(l, e)
			}

			if len(l) >= 3 {
				// do something
				l = l[0:0]
			}
		case <-ctx.Done():
			for v := range ch {
				if e, ok := v.(*entity.Example); ok {
					l = append(l, e)
				}
			}
			if len(l) > 0 {
				// do something
				l = l[0:0]
			}
			return nil
		}
	}
}
