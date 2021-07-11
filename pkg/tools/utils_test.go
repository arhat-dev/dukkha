package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeparateBaseAndImpl(t *testing.T) {
	forBase, forImpl := separateBaseAndImpl("base.", []string{"base.foo", "bar"})

	assert.EqualValues(t, []string{"foo"}, forBase)
	assert.EqualValues(t, []string{"bar"}, forImpl)
}
