package dukkha

import (
	"context"
	"fmt"
	"os"
	"sync"

	"arhat.dev/pkg/exechelper"

	"arhat.dev/dukkha/pkg/types"
	"arhat.dev/dukkha/pkg/utils"
)

type ExecSpecGetFunc func(toExec []string, isFilePath bool) (env, cmd []string, err error)

type ConfigResolvingContext interface {
	Context

	ShellManager
	ToolManager
	TaskManager
	RendererManager
}

// TODO: separate task specific functions to make it a standalone type
type TaskExecContext = Context

type Context interface {
	context.Context

	// rendering
	types.RenderingContext

	ToolUser
	TaskUser

	// DeriveNew Context from this Context
	DeriveNew() Context

	GetBootstrapExecSpec(toExec []string, isFilePath bool) (env, cmd []string, err error)

	Cancel()

	RunTask(ToolKind, ToolName, TaskKind, TaskName) error
	RunShell(shell, script string, isFilePath bool) error

	// dukkha application settings

	FailFast() bool
	ClaimWorkers(n int) int

	ExecValues
}

var (
	_ ConfigResolvingContext = (*dukkhaContext)(nil)

	_ Context = (*dukkhaContext)(nil)
)

// Context of dukkha app, contains global settings and values
type dukkhaContext struct {
	*contextStd

	cache            *sync.Map
	bootstrapExecGen ExecSpecGetFunc

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
	failFast bool
	workers  int
}

func NewConfigResolvingContext(
	parent context.Context,
	globalEnv map[string]string,
	bootstrapExecGen ExecSpecGetFunc,
	failFast bool,
	workers int,
) ConfigResolvingContext {
	ctxStd := newContextStd(parent)
	dukkhaCtx := &dukkhaContext{
		contextStd: ctxStd,

		bootstrapExecGen: bootstrapExecGen,

		contextShells: newContextShells(),
		contextTools:  newContextTools(),
		contextTasks:  newContextTasks(),
		contextExec:   newContextExec(),

		contextRendering: newContextRendering(ctxStd.ctx, globalEnv),

		failFast: failFast,
		workers:  workers,
	}

	return dukkhaCtx
}

func (c *dukkhaContext) GetBootstrapExecSpec(toExec []string, isFilePath bool) (env, cmd []string, err error) {
	return c.bootstrapExecGen(toExec, isFilePath)
}

func (c *dukkhaContext) DeriveNew() Context {
	ctxStd := newContextStd(c.contextStd.ctx)
	newCtx := &dukkhaContext{
		contextStd:       ctxStd,
		cache:            c.cache,
		bootstrapExecGen: c.bootstrapExecGen,

		contextShells:    c.contextShells,
		contextTools:     c.contextTools,
		contextTasks:     c.contextTasks,
		contextRendering: c.contextRendering.clone(ctxStd.ctx),

		// initialized later
		contextExec: nil,

		failFast: c.failFast,
		workers:  c.workers,
	}

	if c.contextExec != nil {
		newCtx.contextExec = c.contextExec.deriveNew()
	} else {
		newCtx.contextExec = newContextExec()
	}

	return newCtx
}

func (c *dukkhaContext) RunTask(k ToolKind, n ToolName, tK TaskKind, tN TaskName) error {
	tool, ok := c.GetTool(k, n)
	if !ok {
		return fmt.Errorf("tool %q with name %q not found", k, n)
	}

	if len(tK) == 0 {
		return fmt.Errorf("invalid empty task kind")
	}

	if len(tN) == 0 {
		return fmt.Errorf("invalid empty task name")
	}

	c.contextExec.thisTool = tool
	c.contextExec.setTask(k, n, tK, tN)

	return c.contextExec.thisTool.Run(c)
}

func (c *dukkhaContext) RunShell(shell, script string, isFilePath bool) error {
	sh, ok := c.allShells[ShellKey{shellName: shell}]
	if !ok {
		return fmt.Errorf("shell %q not found", shell)
	}

	env, cmd, err := sh.GetExecSpec([]string{script}, isFilePath)
	if err != nil {
		return err
	}

	c.AddEnv(env...)

	p, err := exechelper.Do(exechelper.Spec{
		Context: c,
		Command: cmd,
		Env:     c.Env(),

		Stdin: os.Stdin,
		Stderr: utils.PrefixWriter(
			c.OutputPrefix(), c.PrefixColor(),
			c.OutputColor(), os.Stderr,
		),
		Stdout: utils.PrefixWriter(
			c.OutputPrefix(), c.PrefixColor(),
			c.OutputColor(), os.Stdout,
		),
	})
	if err != nil {
		return fmt.Errorf("failed to run script: %w", err)
	}

	code, err := p.Wait()
	if err != nil {
		return fmt.Errorf("command exited with code %d: %w", code, err)
	}

	return nil
}

func (c *dukkhaContext) FailFast() bool {
	return c.failFast
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
