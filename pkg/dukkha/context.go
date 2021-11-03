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

	ExecValues
}

// Context for user facing tasks
type Context interface {
	TaskExecContext

	SetRuntimeOptions(opts RuntimeOptions)
	RunTask(ToolKey, TaskKey) error
}

var (
	_ ConfigResolvingContext = (*dukkhaContext)(nil)

	_ Context = (*dukkhaContext)(nil)
)

// Context of dukkha app, contains global settings and values
type dukkhaContext struct {
	contextStd

	cache *sync.Map

	// shells
	contextShells

	// tools
	contextTools

	// tasks
	contextTasks

	// rendering
	contextRendering

	// task execution, can be null if not running any task
	contextExec
}

func NewConfigResolvingContext(
	parent context.Context,
	ifaceTypeHandler rs.InterfaceTypeHandler,
	globalEnv map[string]string,
) ConfigResolvingContext {
	ctxStd := newContextStd(parent)
	dukkhaCtx := &dukkhaContext{
		contextStd: *ctxStd,

		contextShells: *newContextShells(),
		contextTools:  *newContextTools(),
		contextTasks:  *newContextTasks(),
		contextExec:   *newContextExec(),

		contextRendering: *newContextRendering(
			ctxStd.ctx, ifaceTypeHandler, globalEnv,
		),
	}

	return dukkhaCtx
}

func (c *dukkhaContext) DeriveNew() Context {
	ctxStd := newContextStd(c.contextStd.ctx)
	newCtx := &dukkhaContext{
		contextStd: *ctxStd,
		cache:      c.cache,

		contextShells:    c.contextShells,
		contextTools:     c.contextTools,
		contextTasks:     c.contextTasks,
		contextRendering: *c.contextRendering.clone(ctxStd.ctx),
		contextExec:      *c.contextExec.deriveNew(),
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
