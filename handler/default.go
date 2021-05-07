package handler

import (
	"context"
	"encoding/json"

	"github.com/GGXXLL/core-process/entity"
	"github.com/GGXXLL/core-process/internal/process"
	"github.com/go-kit/kit/log"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type defaultHandler struct {
	logger log.Logger
	db     *gorm.DB
}

func newDefaultHandler(logger log.Logger, db *gorm.DB) process.Out {
	return process.Out{Hs: []process.H{
		&defaultHandler{logger, db},
	}}
}

func (h *defaultHandler) Info() *process.Info {
	return &process.Info{
		Name:      "default",
		BatchSize: 3,
	}
}

func (h *defaultHandler) Handler(ctx context.Context, msg *kafka.Message) (interface{}, error) {
	e := entity.Example{}
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (h *defaultHandler) Batch(ctx context.Context, data []interface{}) error {
	// batch insert to mysql
	rdata := make([]*entity.Example, len(data))
	for i, e := range data {
		rdata[i] = e.(*entity.Example)
	}
	if err := h.db.Table("example").Save(rdata).Error; err != nil {
		return err
	}
	return nil
}
