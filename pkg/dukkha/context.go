package dukkha

import (
	"context"
	"fmt"
	"sync"

	"arhat.dev/rs"
)

type ExecSpecGetFunc func(toExec []string, isFilePath bool) (env Env, cmd []string, err error)

type ConfigResolvingContext interface {
	Context

	ShellManager
	ToolManager
	TaskManager
	RendererManager
}

type TaskExecContext interface {
	context.Context

	RenderingContext
	ShellUser
	ToolUser
	TaskUser

	DeriveNew() Context

	Cancel()

	TranslateANSIStream() bool
	RetainANSIStyle() bool
	ColorOutput() bool
	FailFast() bool
	ClaimWorkers(n int) int

	ExecValues
}

// Context for user facing tasks
type Context interface {
	TaskExecContext

	RunTask(ToolKey, TaskKey) error
}

var (
	_ ConfigResolvingContext = (*dukkhaContext)(nil)

	_ Context = (*dukkhaContext)(nil)
)

// Context of dukkha app, contains global settings and values
type dukkhaContext struct {
	*contextStd

	cache *sync.Map

	// shells
	*contextShells

	// tools
	*contextTools

	// tasks
	*contextTasks

	// rendering
	*contextRendering

	// task execution, can be null if not running any task
	*contextExec

	// application settings
	failFast            bool
	colorOutput         bool
	translateANSIStream bool
	retainANSIStyle     bool
	workers             int
}

func NewConfigResolvingContext(
	parent context.Context,
	ifaceTypeHandler rs.InterfaceTypeHandler,
	failFast bool,
	colorOutput bool,
	translateANSIStream bool,
	retainANSIStyle bool,
	workers int,
	globalEnv map[string]string,
) ConfigResolvingContext {
	ctxStd := newContextStd(parent)
	dukkhaCtx := &dukkhaContext{
		contextStd: ctxStd,

		contextShells: newContextShells(),
		contextTools:  newContextTools(),
		contextTasks:  newContextTasks(),
		contextExec:   newContextExec(),

		contextRendering: newContextRendering(
			ctxStd.ctx, ifaceTypeHandler, globalEnv,
		),

		failFast:            failFast,
		colorOutput:         colorOutput,
		translateANSIStream: translateANSIStream,
		retainANSIStyle:     retainANSIStyle,
		workers:             workers,
	}

	return dukkhaCtx
}

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

func (c *dukkhaContext) AddCache(key, value string) {
	// TBD: runtime cache https://github.com/arhat-dev/dukkha/issues/37
}
