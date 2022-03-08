package handler

import (
	"context"
	"encoding/json"
	"fmt"

	processor "github.com/DoNewsCode/core-processor"
	"github.com/ggxxll/core-processor-demo/entity"
	"github.com/segmentio/kafka-go"
)

// fooHandler implement processor.Handler
type fooHandler struct {
}

func newFooHandler() processor.Out {
	return processor.NewOut(
		&fooHandler{},
	)
}

// Info define the topic name and some config
func (h *fooHandler) Info() *processor.Info {
	return &processor.Info{
		Name:      "foo",
		BatchSize: 2,
	}
}

// Handle decode the *kafka.Message and filter
func (h *fooHandler) Handle(ctx context.Context, msg *kafka.Message) (interface{}, error) {
	// decode
	e := entity.Example{}
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		return nil, err
	}

	// discard unwanted data
	if e.Name == "" {
		return nil, fmt.Errorf("dirty data")
	}
	// If unexpected data is encountered, an exit error is thrown
	if e.Name == "123" {
		return nil, processor.NewFatalErr(fmt.Errorf("error data"))
	}

	// do something

	return &e, nil
}
