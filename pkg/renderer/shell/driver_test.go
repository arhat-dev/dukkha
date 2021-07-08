package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDriver(t *testing.T) {
	d := NewDefault(func(toExec []string, isFilePath bool) (env []string, cmd []string, err error) {
		return
	})

	assert.NotNil(t, d)
}
