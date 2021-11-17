package dukkha

import (
	"context"
	"fmt"
	"path/filepath"

	"arhat.dev/dukkha/pkg/utils"
	"arhat.dev/rs"
	"github.com/huandu/xstrings"
)

type ExecSpecGetFunc func(toExec []string, isFilePath bool) (env Env, cmd []string, err error)

type ConfigResolvingContext interface {
	Context

	ShellManager
	ToolManager
	TaskManager
	RendererManager

	// values
	RendererCacheDir(name string) string
	SetCacheDir(dir string)
	OverrideDefaultGitBranch(branch string)
}

type TaskExecContext interface {
	context.Context

	RenderingContext
	ShellUser
	ToolUser
	TaskUser

	DeriveNew() Context
	Cancel()

	// WithCustomParent divert from current context.Context
	// intended to be only used for defered `after` hooks
	WithCustomParent(parent context.Context) TaskExecContext

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
	*contextStd

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
	globalEnv map[string]utils.LazyValue,
) ConfigResolvingContext {
	ctxStd := newContextStd(parent)
	dukkhaCtx := &dukkhaContext{
		contextStd: ctxStd,

		contextShells: *newContextShells(),
		contextTools:  *newContextTools(),
		contextTasks:  *newContextTasks(),
		contextExec:   *newContextExec(),

		contextRendering: *newContextRendering(
			ctxStd, ifaceTypeHandler, globalEnv,
		),
	}

	return dukkhaCtx
}

func (c *dukkhaContext) DeriveNew() Context {
	return c.deriveNew(c.contextStd.ctx, true)
}

func (c *dukkhaContext) deriveNew(parent context.Context, deepCopy bool) Context {
	ctxStd := newContextStd(parent)
	newCtx := &dukkhaContext{
		contextStd: ctxStd,

		contextShells:    c.contextShells,
		contextTools:     c.contextTools,
		contextTasks:     c.contextTasks,
		contextRendering: *c.contextRendering.clone(ctxStd, deepCopy),
		contextExec:      *c.contextExec.deriveNew(),
	}

	return newCtx
}

func (c *dukkhaContext) RendererCacheDir(name string) string {
	return filepath.Join(c.CacheDir(), "renderer", xstrings.ToKebabCase(name))
}

func (c *dukkhaContext) RunTask(k ToolKey, tK TaskKey) error {
	tool, ok := c.GetTool(k)
	if !ok {
		return fmt.Errorf("tool %q not found", k)
	}

	return tool.Run(c, tK)
}

func (c *dukkhaContext) WithCustomParent(parent context.Context) TaskExecContext {
	return c.deriveNew(parent, false)
}
