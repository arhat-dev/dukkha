package dukkha

import (
	"arhat.dev/dukkha/pkg/field"
)

type Shell interface {
	field.Field

	GetExecSpec(toExec []string, isFilePath bool) (env, cmd []string, err error)
}

type ShellUser interface {
	GetShell(name string) (Shell, bool)
	AllShells() map[string]Shell
}

type ShellManager interface {
	AddShell(name string, impl Shell)
}

func newContextShells() *contextShells {
	return &contextShells{
		allShells: make(map[string]Shell),
	}
}

type contextShells struct {
	allShells map[string]Shell
}

func (c *contextShells) GetShell(name string) (Shell, bool) {
	sh, ok := c.allShells[name]
	return sh, ok
}

func (c *contextShells) AddShell(name string, impl Shell) {
	c.allShells[name] = impl
}

func (c *contextShells) AllShells() map[string]Shell {
	return c.allShells
}
