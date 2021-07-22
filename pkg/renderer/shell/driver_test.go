package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDriver(t *testing.T) {
	d := NewDefault()

	assert.NotNil(t, d)
}
