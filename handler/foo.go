package handler

import (
	"context"
	"encoding/json"

	"github.com/DoNewsCode/core/otkafka/processor"
	"github.com/GGXXLL/core-processor-demo/entity"
	"github.com/go-kit/kit/log"
	"github.com/segmentio/kafka-go"
)

// fooHandler implement processor.Handler
type fooHandler struct {
	logger log.Logger
}

func newFooHandler(logger log.Logger) processor.Out {
	return processor.NewOut(
		&fooHandler{logger: logger},
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

	// filter dirty data
	if e.Name == "" {
		return nil, nil
	}
	// do something

	return &e, nil
}
