package handler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ggxxll/core-processor-demo/entity"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func Test_defaultHandler_Handle(t *testing.T) {

	data := &entity.Example{
		Name: "test",
	}
	b, _ := json.Marshal(data)

	handler := defaultHandler{}
	res, err := handler.Handle(context.Background(), &kafka.Message{Value: b})
	assert.NoError(t, err)

	assert.Equal(t, data, res.(*entity.Example))
}

func Test_defaultHandler_Batch(t *testing.T) {
	db, closer := setupDB()
	defer closer()
	handler := defaultHandler{db: db}
	data := []interface{}{
		&entity.Example{
			Name: "test1",
		},
		&entity.Example{
			Name: "test2",
		},
	}
	err := handler.Batch(context.Background(), data)
	assert.NoError(t, err)
	var count int64 = 0
	err = db.Model(&entity.Example{}).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(len(data)), count)
}
