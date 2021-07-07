package dukkha

type Shell interface {
	GetExecSpec(toExec []string, isFilePath bool) (env, cmd []string, err error)
}

type ShellManager interface {
	AddShell(name string, impl Shell)
}

type ShellKey struct {
	shellName string
}

func newContextShells() *contextShells {
	return &contextShells{
		allShells: make(map[ShellKey]Shell),
	}
}

type contextShells struct {
	allShells map[ShellKey]Shell
}

func (c *contextShells) AddShell(name string, impl Shell) {
	c.allShells[ShellKey{shellName: name}] = impl
}
