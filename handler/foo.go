package handler

import (
	"context"
	"encoding/json"
	"github.com/GGXXLL/core-kafka/entity"
	"github.com/GGXXLL/core-kafka/internal/process"
	"github.com/go-kit/kit/log"
)

type fooHandler struct {
	logger log.Logger
}

func newFooHandler(logger log.Logger) process.Out {
	return process.Out{Hs: []process.H{&fooHandler{logger}}}
}

func (h *fooHandler) Name() string {
	return "foo"
}

func (h *fooHandler) Handler(ctx context.Context, msg *kafka.Message) (interface{}, error) {
	e := entity.Example{}
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (h *fooHandler) Batch(ctx context.Context, ch chan interface{}) error {
	<-ctx.Done()
	return nil
}
