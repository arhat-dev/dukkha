package shell

import (
	"testing"

	"arhat.dev/dukkha/pkg/dukkha"
	"github.com/stretchr/testify/assert"
)

var _ dukkha.Renderer = (*Driver)(nil)

func TestNewDriver(t *testing.T) {
	d := NewDefault("")

	assert.NotNil(t, d)
}
