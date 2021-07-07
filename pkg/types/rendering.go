package types

import "context"

type RenderingContext interface {
	context.Context

	ImmutableValues
	MutableValues

	Env() map[string]string

	RenderYaml(renderer string, rawData interface{}) (result string, err error)
}

type ImmutableValues interface {
	WorkingDir() string
	CacheDir() string

	GitBranch() string
	GitWorkTreeClean() bool
	GitTag() string
	GitDefaultBranch() string
	GitCommit() string

	HostKernel() string
	HostKernelVersion() string
	HostArch() string
	HostOS() string
	HostOSVersion() string
}

type MutableValues interface {
	SetMatrixFilter(mf map[string][]string)
	MatrixFilter() map[string][]string

	MatrixArch() string
	MatrixKernel() string
	MatrixLibc() string

	AddEnv(env ...string)
}
