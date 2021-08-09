package dukkha

import (
	"context"
	"time"
)

var (
	_ context.Context = (*contextStd)(nil)
)

func newContextStd(parent context.Context) *contextStd {
	stdCtx := &contextStd{}
	stdCtx.ctx, stdCtx.cancel = context.WithCancel(parent)
	return stdCtx
}

type contextStd struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *contextStd) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *contextStd) Err() error {
	return c.ctx.Err()
}

func (c *contextStd) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c *contextStd) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func (c *contextStd) Cancel() {
	c.cancel()
}
