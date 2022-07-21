package dukkha

import (
	"context"
	"fmt"
	"path"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/pathhelper"
	"arhat.dev/rs"
	"arhat.dev/tlang"
	"github.com/huandu/xstrings"
)

type ExecSpecGetFunc func(toExec []string, isFilePath bool) (env Env, cmd []string, err error)

// ConfigResolvingContext interface definition for config resolving
type ConfigResolvingContext interface {
	Context

	ShellManager
	ToolManager
	TaskManager
	RendererManager

	// cache fs
	ToolCacheFS(t Tool) *fshelper.OSFS
	TaskCacheFS(t Task) *fshelper.OSFS
	RendererCacheFS(name string) *fshelper.OSFS
}

// TaskExecContext interface definition for task execution
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

// Context of dukkha app, contains global settings and values
type dukkhaContext struct {
	// rendering
	// MUST be the first element, we are casting its pointer to get pointer to dukkhaContext
	contextRendering

	// shells
	contextShells

	// tools
	contextTools

	// tasks
	contextTasks

	// task execution, can be null if not running any task
	contextExec
}

func NewConfigResolvingContext(
	parent context.Context,
	ifaceTypeHandler rs.InterfaceTypeHandler,
	globalEnv map[string]tlang.LazyValueType[string],
) ConfigResolvingContext {
	ctxStd := newContextStd(parent)
	dukkhaCtx := &dukkhaContext{
		contextShells: newContextShells(),
		contextTools:  newContextTools(),
		contextTasks:  newContextTasks(),
		contextExec:   newContextExec(),

		contextRendering: newContextRendering(
			ctxStd, ifaceTypeHandler, globalEnv,
		),
	}

	return dukkhaCtx
}

func (c *dukkhaContext) DeriveNew() Context {
	return c.deriveNew(c.contextStd.ctx, true)
}

func (c *dukkhaContext) deriveNew(parent context.Context, deepCopy bool) Context {
	newCtx := &dukkhaContext{
		contextRendering: c.contextRendering.clone(newContextStd(parent), deepCopy),
		contextShells:    c.contextShells,
		contextTools:     c.contextTools,
		contextTasks:     c.contextTasks,
		contextExec:      c.contextExec.deriveNew(),
	}

	return newCtx
}

func replaceInvalidWindowsPathChars(name string) string {
	// replace invalid characters for windows
	basename := []rune(xstrings.ToKebabCase(name))
	for i, c := range basename {
		if pathhelper.IsReservedWindowsPathChar(c) {
			basename[i] = '-'
		}
	}

	return string(basename)
}

func (c *dukkhaContext) RendererCacheFS(name string) *fshelper.OSFS {
	name = replaceInvalidWindowsPathChars(name)
	return lazyEnsuredSubFS(c.cacheFS, path.Join("renderer", name))
}

func (c *dukkhaContext) ToolCacheFS(t Tool) *fshelper.OSFS {
	k := string(t.Kind())
	if len(k) == 0 {
		panic("invalid empty tool kind")
	}

	name := replaceInvalidWindowsPathChars(string(t.Name()))
	if len(name) == 0 {
		name = "_"
	}

	return lazyEnsuredSubFS(c.cacheFS, path.Join(k, name))
}

func (c *dukkhaContext) TaskCacheFS(t Task) *fshelper.OSFS {
	toolKind := string(t.ToolKind())
	if len(toolKind) == 0 {
		panic("invalid empty tool kind")
	}

	toolName := replaceInvalidWindowsPathChars(string(t.ToolName()))
	if len(toolName) == 0 {
		toolName = "_"
	}

	kind := string(t.Kind())
	if len(kind) == 0 {
		panic("invalid empty task kind")
	}

	name := replaceInvalidWindowsPathChars(string(t.Name()))
	if len(name) == 0 {
		panic("invalid empty task name")
	}

	return lazyEnsuredSubFS(c.cacheFS, path.Join(toolKind, toolName, kind, name))
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
