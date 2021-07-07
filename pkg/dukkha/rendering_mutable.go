package dukkha

import (
	"strings"
	"sync"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/types"
)

func newContextMutableValues() *mutableValues {
	return &mutableValues{
		matrixFilter: nil,
		env:          make(map[string]string),
		mu:           new(sync.RWMutex),
	}
}

var _ types.MutableValues = (*mutableValues)(nil)

type mutableValues struct {
	matrixFilter map[string][]string

	env map[string]string
	mu  *sync.RWMutex
}

func (c *mutableValues) clone() *mutableValues {
	newValues := newContextMutableValues()

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.matrixFilter) != 0 {
		c.matrixFilter = make(map[string][]string)
		for k, v := range c.matrixFilter {
			c.matrixFilter[k] = sliceutils.NewStrings(v)
		}
	}

	for k, v := range c.env {
		newValues.env[k] = v
	}

	return newValues
}

func (c *mutableValues) SetMatrixFilter(mf map[string][]string) {
	c.matrixFilter = mf
}

func (c *mutableValues) MatrixFilter() map[string][]string {
	return c.matrixFilter
}

func (c *mutableValues) MatrixArch() string {
	return c.env[constant.ENV_MATRIX_ARCH]
}

func (c *mutableValues) MatrixKernel() string {
	return c.env[constant.ENV_MATRIX_KERNEL]
}

func (c *mutableValues) MatrixLibc() string {
	return c.env[constant.ENV_MATRIX_LIBC]
}

func (c *mutableValues) AddEnv(entries ...string) {
	for _, entry := range entries {
		parts := strings.SplitN(entry, "=", 2)
		key, value := parts[0], ""
		if len(parts) == 2 {
			value = parts[1]
		}

		c.env[key] = value
	}
}
