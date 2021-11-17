package dukkha

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/utils"
)

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
		matchFilterStr := strings.TrimRight(strings.TrimSpace(callArgs[1]), ",")
		mf := make(map[string][]string)
		err = yaml.Unmarshal([]byte(matchFilterStr), &mf)
		if err != nil {
			return nil, fmt.Errorf("invalid matrix arg %q: %w", callArgs[1], err)
		}

		ref.MatrixFilter = matrix.NewFilter(mf)
	}

	return ref, nil
}

type TaskReference struct {
	ToolKind ToolKind
	ToolName ToolName
	TaskKind TaskKind
	TaskName TaskName

	MatrixFilter *matrix.Filter
}

func (t *TaskReference) ToolKey() ToolKey {
	return ToolKey{Kind: t.ToolKind, Name: t.ToolName}
}

func (t *TaskReference) TaskKey() TaskKey {
	return TaskKey{Kind: t.TaskKind, Name: t.TaskName}
}
