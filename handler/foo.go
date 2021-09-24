package handler

import (
	"context"
	"encoding/json"

	processor "github.com/DoNewsCode/core-processor"
	"github.com/DoNewsCode/core/logging"

	"github.com/ggxxll/core-processor-demo/entity"
	"github.com/go-kit/kit/log"
	"github.com/segmentio/kafka-go"
)

// fooHandler implement processor.Handler
type fooHandler struct {
	logger logging.LevelLogger
}

func newFooHandler(logger log.Logger) processor.Out {
	return processor.NewOut(
		&fooHandler{logger: logging.WithLevel(logger)},
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
		h.logger.Err(err)
		return nil, err
	}

	// filter dirty data
	if e.Name == "" {
		return nil, nil
	}
	// do something

	return &e, nil
}
