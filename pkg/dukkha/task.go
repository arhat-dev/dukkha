package dukkha

import (
	"arhat.dev/pkg/fshelper"

	"arhat.dev/dukkha/pkg/matrix"
)

type (
	TaskKind string
	TaskName string
)

type TaskKey struct {
	Kind TaskKind
	Name TaskName
}

func (k TaskKey) String() string { return string(k.Kind) + ":" + string(k.Name) }

// Task implementation requirements
type Task interface {
	Resolvable

	// ToolKind is the tool managing this task (e.g. docker)
	ToolKind() ToolKind

	// ToolName it the name of the tool this task belongs to
	ToolName() ToolName

	// Kind of the task (e.g. build)
	Kind() TaskKind

	// Name of the task (e.g. foo)
	Name() TaskName

	// Key is the composite key of task kind and name
	Key() TaskKey

	// Init this task
	Init(cacheFS *fshelper.OSFS) error

	// GetMatrixSpecs for matrix execution
	//
	// The implementation MUST be thread safe
	GetMatrixSpecs(rc RenderingContext) ([]matrix.Entry, error)

	// GetExecSpecs generate commands using current field values
	//
	// The implementation MUST be thread safe
	GetExecSpecs(rc TaskExecContext, options TaskMatrixExecOptions) ([]TaskExecSpec, error)

	// GetHookExecSpecs generate hook run target
	//
	// The implementation MUST be thread safe
	GetHookExecSpecs(rc TaskExecContext, state TaskExecStage) ([]TaskExecSpec, error)
}
