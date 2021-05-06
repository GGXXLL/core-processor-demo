package handler

import (
	"context"
	"encoding/json"

	"github.com/GGXXLL/core-process/entity"
	"github.com/GGXXLL/core-process/internal/process"
	"github.com/go-kit/kit/log"
	"github.com/segmentio/kafka-go"
)

type fooHandler struct {
	logger log.Logger
}

func newFooHandler(logger log.Logger) process.Out {
	return process.Out{Hs: []process.H{
		&fooHandler{logger},
	}}
}

// Info define the topic name and some config
func (h *fooHandler) Info() *process.Info {
	return &process.Info{
		Name:      "foo",
		BatchSize: 2,
	}
}

// Handler decode the *kafka.Message and filter
func (h *fooHandler) Handler(ctx context.Context, msg *kafka.Message) (interface{}, error) {
	// decode
	e := entity.Example{}
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		return nil, err
	}

	// filter
	if e.Id <= -1 {
		return nil, nil
	}

	return &e, nil
}

// Batch this func does nothing, just implement process.H
func (h *fooHandler) Batch(ctx context.Context, data []interface{}) error {
	return nil
}
