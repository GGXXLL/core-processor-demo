package handler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ggxxll/core-processor-demo/entity"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func Test_fooHandler_Handle(t *testing.T) {
	handler := fooHandler{}
	cases := []struct {
		name   string
		data   *entity.Example
		assert func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{name: "correct", data: &entity.Example{Name: "test"}, assert: assert.NotNil},
		{name: "dirty", data: &entity.Example{Name: ""}, assert: assert.Nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, _ := json.Marshal(c.data)
			res, err := handler.Handle(context.Background(), &kafka.Message{Value: b})
			assert.NoError(t, err)

			c.assert(t, res)
		})
	}
}
