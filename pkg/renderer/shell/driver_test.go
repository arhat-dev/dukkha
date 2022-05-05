package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
)

var _ dukkha.Renderer = (*Driver)(nil)

func TestNewDriver(t *testing.T) {
	t.Parallel()

	d := NewDefault("")

	assert.NotNil(t, d)
}
