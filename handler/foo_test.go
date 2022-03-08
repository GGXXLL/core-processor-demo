package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	processor "github.com/DoNewsCode/core-processor"
	"github.com/ggxxll/core-processor-demo/entity"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func Test_fooHandler_Handle(t *testing.T) {
	handler := fooHandler{}
	cases := []struct {
		name    string
		data    *entity.Example
		assert  func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
		wantErr error
	}{
		{name: "correct", data: &entity.Example{Name: "test"}, assert: assert.NotNil, wantErr: nil},
		{name: "dirty", data: &entity.Example{Name: ""}, assert: assert.Nil, wantErr: errors.New("dirty data")},
		{name: "unexpect", data: &entity.Example{Name: "123"}, assert: assert.Nil, wantErr: processor.NewFatalErr(fmt.Errorf("error data"))},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, _ := json.Marshal(c.data)
			res, err := handler.Handle(context.Background(), &kafka.Message{Value: b})

			if err != nil && errors.Is(err, c.wantErr) {
				t.Fatalf("want err %v, got %v", c.wantErr, err)
			}

			c.assert(t, res)
		})
	}
}
