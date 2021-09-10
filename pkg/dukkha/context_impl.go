package dukkha

import "fmt"

func (c *dukkhaContext) DeriveNew() Context {
	ctxStd := newContextStd(c.contextStd.ctx)
	newCtx := &dukkhaContext{
		contextStd: ctxStd,
		cache:      c.cache,

		contextShells:    c.contextShells,
		contextTools:     c.contextTools,
		contextTasks:     c.contextTasks,
		contextRendering: c.contextRendering.clone(ctxStd.ctx),

		// initialized later
		contextExec: nil,

		failFast:            c.failFast,
		colorOutput:         c.colorOutput,
		translateANSIStream: c.translateANSIStream,
		retainANSIStyle:     c.retainANSIStyle,
		workers:             c.workers,
	}

	if c.contextExec != nil {
		newCtx.contextExec = c.contextExec.deriveNew()
	} else {
		newCtx.contextExec = newContextExec()
	}

	return newCtx
}

func (c *dukkhaContext) RunTask(k ToolKey, tK TaskKey) error {
	tool, ok := c.GetTool(k)
	if !ok {
		return fmt.Errorf("tool %q not found", k)
	}

	c.contextExec.SetTask(k, tK)
	return tool.Run(c)
}

func (c *dukkhaContext) FailFast() bool {
	return c.failFast
}

func (c *dukkhaContext) ColorOutput() bool {
	return c.colorOutput
}

func (c *dukkhaContext) TranslateANSIStream() bool {
	return c.translateANSIStream
}

func (c *dukkhaContext) RetainANSIStyle() bool {
	return c.retainANSIStyle
}

func (c *dukkhaContext) ClaimWorkers(n int) int {
	if c.workers > n {
		return n
	}

	// TODO: limit workers
	return c.workers
}
