package handler

import (
	"context"
	"encoding/json"

	processor "github.com/DoNewsCode/core-processor"
	"github.com/DoNewsCode/core/contract"
	"github.com/ggxxll/core-processor-demo/entity"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

// defaultHandler implement processor.BatchHandler
type defaultHandler struct {
	db  *gorm.DB
	env contract.Env
}

func newDefaultHandler(db *gorm.DB, env contract.Env) processor.Out {
	return processor.NewOut(
		&defaultHandler{db, env},
	)
}

func (h *defaultHandler) Info() *processor.Info {
	// using env to distinguish the specific configuration
	if h.env.IsLocal() {
		return &processor.Info{
			Name:      "default",
			BatchSize: 1,
		}
	}
	return &processor.Info{
		Name:      "default",
		BatchSize: 10,
	}
}

func (h *defaultHandler) Handle(ctx context.Context, msg *kafka.Message) (interface{}, error) {
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
	if err := h.db.Table("examples").Save(rdata).Error; err != nil {
		return err
	}
	return nil
}
