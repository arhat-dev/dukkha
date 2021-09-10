package dukkha

import (
	"context"
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

type ContextOptions struct {
	InterfaceTypeHandler rs.InterfaceTypeHandler
	FailFast             bool
	ColorOutput          bool
	TranslateANSIStream  bool
	RetainANSIStyle      bool
	Workers              int
	GlobalEnv            map[string]string
}

func NewConfigResolvingContext(
	parent context.Context,
	opts ContextOptions,
) ConfigResolvingContext {
	ctxStd := newContextStd(parent)
	dukkhaCtx := &dukkhaContext{
		contextStd: ctxStd,

		contextShells: newContextShells(),
		contextTools:  newContextTools(),
		contextTasks:  newContextTasks(),
		contextExec:   newContextExec(),

		contextRendering: newContextRendering(
			ctxStd.ctx, opts.InterfaceTypeHandler, opts.GlobalEnv,
		),

		failFast:            opts.FailFast,
		colorOutput:         opts.ColorOutput,
		translateANSIStream: opts.TranslateANSIStream,
		retainANSIStyle:     opts.RetainANSIStyle,
		workers:             opts.Workers,
	}

	return dukkhaCtx
}
