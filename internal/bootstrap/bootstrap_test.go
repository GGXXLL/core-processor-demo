package bootstrap

import (
	"os"
	"testing"
)

func TestBootstrap(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Error occurred: [%s]", r)
		}
	}()

	os.Args = []string{"", "--config", "../../config.yaml"}
	_, shutdown := Bootstrap()
	shutdown()
}
