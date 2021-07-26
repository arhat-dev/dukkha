package dukkha

import (
	"fmt"
	"io"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/utils"
)

type (
	TaskKind string
	TaskName string

	TaskKey struct {
		Kind TaskKind
		Name TaskName
	}
)

func (k TaskKey) String() string {
	return string(k.Kind) + ":" + string(k.Name)
}

type TaskReference struct {
	ToolKind ToolKind
	ToolName ToolName
	TaskKind TaskKind
	TaskName TaskName

	MatrixFilter map[string][]string
}

func (t *TaskReference) ToolKey() ToolKey {
	return ToolKey{Kind: t.ToolKind, Name: t.ToolName}
}

func (t *TaskReference) TaskKey() TaskKey {
	return TaskKey{Kind: t.TaskKind, Name: t.TaskName}
}

// ParseTaskReference parse task ref
//
// <tool-kind>{:<tool-name>}:<task-kind>(<task-name>, ...)
//
// e.g. buildah:build(dukkha) # use default matrix
// 		buildah:build(dukkha, {kernel: [linux]}) # use custom matrix
//		buildah:in-docker:build(dukkha, {kernel: [linux]}) # with tool-name
func ParseTaskReference(taskRef string, defaultToolName ToolName) (*TaskReference, error) {
	callStart := strings.IndexByte(taskRef, '(')
	if callStart < 0 {
		return nil, fmt.Errorf("missing task call `(<task-name>)`")
	}

	ref := &TaskReference{}

	// <tool-kind>{:<tool-name>}:<task-kind>
	parts := strings.Split(taskRef[:callStart], ":")
	ref.ToolKind = ToolKind(parts[0])

	switch len(parts) {
	case 2:
		// no tool name set, use the default tool name
		// no matter what kind the tool is
		//
		// current task
		// 		buildah:in-docker:build 	# tool name is `in-docker`
		// has task reference in hook:
		// 		buildah:login(foo)    	# same kind
		// 		golang:build(bar)		# different kind
		// will actually be treated as
		// 		buildah:in-docker:login(foo)	# same kind
		//		golang:in-docker:build(bar)		# different kind

		ref.ToolName = defaultToolName
		ref.TaskKind = TaskKind(parts[1])
	case 3:
		ref.ToolName = ToolName(parts[1])
		ref.TaskKind = TaskKind(parts[2])
	default:
		return nil, fmt.Errorf("invalid tool reference %q", taskRef)
	}

	call, err := utils.ParseBrackets(taskRef[callStart+1:])
	if err != nil {
		return nil, fmt.Errorf("invalid task call: %w", err)
	}
	callArgs := strings.SplitN(call, ",", 2)
	ref.TaskName = TaskName(strings.TrimSpace(callArgs[0]))

	switch len(callArgs) {
	case 1:
		// using default matrix spec, do nothing
	case 2:
		// second arg is matrix spec
		matrixFilterStr := strings.TrimRight(strings.TrimSpace(callArgs[1]), ",")
		ref.MatrixFilter = make(map[string][]string)
		err = yaml.Unmarshal([]byte(matrixFilterStr), &ref.MatrixFilter)
		if err != nil {
			return nil, fmt.Errorf("invalid matrix arg\n\n%s\nerror: %w", callArgs[1], err)
		}
	}

	return ref, nil
}

type TaskExecStage uint8

const (
	StageBefore TaskExecStage = iota + 1

	StageBeforeMatrix
	StageAfterMatrixSuccess
	StageAfterMatrixFailure
	StageAfterMatrix

	StageAfterSuccess
	StageAfterFailure
	StageAfter
)

func (s TaskExecStage) String() string {
	return map[TaskExecStage]string{
		StageBefore: "before",

		StageBeforeMatrix:       "before:matrix",
		StageAfterMatrixSuccess: "after:matrix:success",
		StageAfterMatrixFailure: "after:matrix:failure",
		StageAfterMatrix:        "after:matrix",

		StageAfterSuccess: "after:success",
		StageAfterFailure: "after:failure",
		StageAfter:        "after",
	}[s]
}

type RunTaskOrRunCmd interface{}

type TaskExecSpec struct {
	// Delay execution
	Delay time.Duration

	// OutputAsReplace to replace same string in following TaskExecSpecs
	OutputAsReplace string

	FixOutputForReplace func(newValue []byte) []byte

	Chdir string

	Env     []string
	Command []string

	AlterExecFunc func(
		replace map[string][]byte,
		stdin io.Reader, stdout, stderr io.Writer,
	) (RunTaskOrRunCmd, error)

	Stdin io.Reader

	IgnoreError bool

	// UseShell if true, write command to local script cache
	// and execute with the target shell (as referenced by ShellName)
	UseShell bool

	// ShellName to reference a shell to execute command
	// when `UseShell` is true
	//
	// the availability of the shells denpends on `shells` in dukkha config
	// a special shell name is `bootstrap`, which will use the bootstrap
	// section as shell interpreter
	ShellName string
}

type Task interface {
	field.Field

	// Kind of the tool managing this task (e.g. docker)
	ToolKind() ToolKind

	// Name of the tool managing this task (e.g. my-tool)
	ToolName() ToolName

	// Kind of the task (e.g. build)
	Kind() TaskKind

	// Name of the task (e.g. foo)
	Name() TaskName

	// Key of this task
	Key() TaskKey

	// GetMatrixSpecs for matrix execution
	//
	// The implementation MUST be safe to be used concurrently
	GetMatrixSpecs(rc RenderingContext) ([]matrix.Entry, error)

	// GetExecSpecs generate commands using current field values
	//
	// The implementation MUST be safe to be used concurrently
	GetExecSpecs(rc TaskExecContext, options TaskMatrixExecOptions) ([]TaskExecSpec, error)

	// GetHookExecSpecs generate hook run target
	//
	// The implementation MUST be safe to be used concurrently
	GetHookExecSpecs(ctx TaskExecContext, state TaskExecStage) ([]RunTaskOrRunCmd, error)

	// DoAfterFieldsResolved is a helper function to ensure no data race
	//
	// The implementation MUST be safe to be used concurrently
	DoAfterFieldsResolved(
		rc RenderingContext,
		depth int,
		do func() error,
		fieldNames ...string,
	) error

	ContinueOnError() bool
}

type TaskManager interface {
	AddToolSpecificTasks(kind ToolKind, name ToolName, tasks []Task)
}

type TaskUser interface {
	GetToolSpecificTasks(k ToolKey) ([]Task, bool)
	AllToolSpecificTasks() map[ToolKey][]Task
}

func newContextTasks() *contextTasks {
	return &contextTasks{
		toolSpecificTasks: make(map[ToolKey][]Task),
	}
}

type contextTasks struct {
	toolSpecificTasks map[ToolKey][]Task
}

func (c *contextTasks) AddToolSpecificTasks(k ToolKind, n ToolName, tasks []Task) {
	toolKey := ToolKey{Kind: k, Name: n}

	c.toolSpecificTasks[toolKey] = append(
		c.toolSpecificTasks[toolKey], tasks...,
	)
}

func (c *contextTasks) GetToolSpecificTasks(k ToolKey) ([]Task, bool) {
	tasks, ok := c.toolSpecificTasks[k]
	return tasks, ok
}

func (c *contextTasks) AllToolSpecificTasks() map[ToolKey][]Task {
	return c.toolSpecificTasks
}
