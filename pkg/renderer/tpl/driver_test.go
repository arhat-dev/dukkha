package tpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDriver(t *testing.T) {
	assert.NotNil(t, NewDefault(""))
}
