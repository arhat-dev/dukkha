package dukkha

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/types"
	"arhat.dev/dukkha/pkg/utils"
)

type TaskReference struct {
	ToolKind ToolKind
	ToolName ToolName
	TaskKind TaskKind
	TaskName TaskName

	MatrixFilter map[string][]string
}

func (r *TaskReference) HasToolName() bool {
	return len(r.ToolName) != 0
}

// ParseTaskReference parse task ref
//
// <tool-kind>{:<tool-name>}:<task-kind>(<task-name>, ...)
//
// e.g. buildah:bud(dukkha) # use default matrix
// 		buildah:bud(dukkha, {kernel: [linux]}) # use custom matrix
//		buildah:in-docker:bud(dukkha, {kernel: [linux]}) # with tool-name
func ParseTaskReference(taskRef string) (*TaskReference, error) {
	callStart := strings.IndexByte(taskRef, '(')
	if callStart < 0 {
		return nil, fmt.Errorf("missing task call `()`")
	}

	call, err := utils.ParseBrackets(taskRef[callStart+1:])
	if err != nil {
		return nil, fmt.Errorf("invalid task call: %w", err)
	}

	ref := &TaskReference{}
	callArgs := strings.SplitN(call, ",", 2)
	ref.TaskName = TaskName(strings.TrimSpace(callArgs[0]))

	switch len(callArgs) {
	case 1:
		// using default matrix spec, do nothing
	case 2:
		matrixFilterStr := strings.TrimRight(strings.TrimSpace(callArgs[1]), ",")
		ref.MatrixFilter = make(map[string][]string)
		err = yaml.Unmarshal([]byte(matrixFilterStr), &ref.MatrixFilter)
		if err != nil {
			return nil, fmt.Errorf("invalid matrix arg \n\n%s\nerror: %w", callArgs[1], err)
		}
	}

	parts := strings.Split(taskRef[:callStart], ":")
	ref.ToolKind = ToolKind(parts[0])

	switch len(parts) {
	case 2:
		ref.TaskKind = TaskKind(parts[1])
	case 3:
		ref.ToolName = ToolName(parts[1])
		ref.TaskKind = TaskKind(parts[2])
	default:
		return nil, fmt.Errorf("invalid prefix %q", taskRef)
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

// TaskType for interface type registration
var TaskType = reflect.TypeOf((*Task)(nil)).Elem()

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
	) ([]TaskExecSpec, error)

	Stdin io.Reader

	IgnoreError bool
}

type Task interface {
	types.Field

	// Kind of the tool managing this task (e.g. docker)
	ToolKind() ToolKind

	// Name of the tool managing this task (e.g. my-tool)
	ToolName() ToolName

	// Kind of the task (e.g. build)
	Kind() TaskKind

	// Name of the task
	Name() TaskName

	// GetMatrixSpecs for matrix build
	GetMatrixSpecs(rc types.RenderingContext) ([]types.MatrixSpec, error)

	// GetExecSpecs generate commands using current field values
	GetExecSpecs(rc types.RenderingContext, toolCmd []string) ([]TaskExecSpec, error)

	RunHooks(taskCtx Context, state TaskExecStage) error
}

type TaskManager interface {
	TaskUser

	AddToolSpecificTasks(kind ToolKind, name ToolName, tasks []Task)
}

type TaskUser interface {
	GetToolSpecificTasks(kind ToolKind, name ToolName) ([]Task, bool)
	AllToolSpecificTasks() map[ToolKey][]Task
}

type (
	TaskKind string
	TaskName string
)

type TaskKey struct {
	Kind TaskKind
	Name TaskName
}

func (k TaskKey) String() string {
	return string(k.Kind) + ":" + string(k.Name)
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

func (c *contextTasks) GetToolSpecificTasks(k ToolKind, n ToolName) ([]Task, bool) {
	tasks, ok := c.toolSpecificTasks[ToolKey{Kind: k, Name: n}]
	return tasks, ok
}

func (c *contextTasks) AllToolSpecificTasks() map[ToolKey][]Task {
	return c.toolSpecificTasks
}
