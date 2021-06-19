package tools

import (
	"reflect"

	"arhat.dev/dukkha/pkg/field"
)

// TaskConfigType for tools.TaskConfig interface type registration
var TaskConfigType = reflect.TypeOf((*TaskConfig)(nil)).Elem()

type TaskConfig interface {
	field.Interface

	// Kind of the task, e.g. golang:build
	Kind() string
}

type BaseTask struct {
	field.BaseField

	Name string `yaml:"name"`
}

type BaseMatrixConfig struct {
	field.BaseField

	Include map[string][]string `yaml:"include"`
	Exclude map[string][]string `yaml:"exclude"`

	OS   []string `yaml:"os"`
	Arch []string `yaml:"arch"`
}

// func resolveTasks() {
// 	// taskParts := strings.Split(taskField.Name, ":")

// 	// var (
// 	// 	toolName = taskParts[0]
// 	// 	toolID   string
// 	// 	taskType string
// 	// )

// 	// switch len(taskParts) {
// 	// case 2:
// 	// 	taskType = taskParts[1]
// 	// case 3:
// 	// 	toolID, taskType = taskParts[1], taskParts[2]
// 	// default:
// 	// 	return fmt.Errorf(
// 	// 		"task: invalid task field %q, expecting 1 or 2 colon, got %d",
// 	// 		taskField.Name, len(taskParts),
// 	// 	)
// 	// }

// 	// key, err := tools.CreateTaskTypeKey(toolName, toolID, taskType)
// 	// if err != nil {
// 	// 	return fmt.Errorf("task: invalid task field: %w", err)
// 	// }

// 	// c.Tasks[*key] = nil
// 	// if len(taskField.Renderer) != 0 {
// 	// 	// requires extra rendering
// 	// 	strVal, ok := data.(string)
// 	// 	if !ok {
// 	// 		return fmt.Errorf("task.%s: unexpected non string value", key.String())
// 	// 	}

// 	// 	// TODO: mark to be processed later
// 	// 	_ = strVal
// 	// } else {
// 	// 	// can unmarshal childs directly
// 	// 	tasksBytes, err := yaml.Marshal(data)
// 	// 	if err != nil {
// 	// 		return fmt.Errorf("task.%s: marhsal: %w", key.String(), err)
// 	// 	}

// 	// 	taskConfigType, err := tools.GetTaskConfigType(toolName, taskType)
// 	// 	if err != nil {
// 	// 		return fmt.Errorf("task.%s", key.String())
// 	// 	}

// 	// 	_ = tasksBytes
// 	// 	_ = taskConfigType
// 	// }
// }
