package tools

import (
	"context"
	"os"
	"reflect"

	"arhat.dev/pkg/exechelper"

	"arhat.dev/dukkha/pkg/field"
)

// ToolType for interface type registration
var ToolType = reflect.TypeOf((*Tool)(nil)).Elem()

// nolint:revive
type Tool interface {
	// Kind of the tool, e.g. golang, docker
	Kind() string

	ResolveTasks(tasks []Task) error
}

type BaseTool struct {
	field.BaseField

	Name string   `yaml:"name"`
	Path string   `yaml:"path"`
	Env  []string `yaml:"env"`

	GlobalArgs []string `yaml:"args"`
}

func (t *BaseTool) Exec(ctx context.Context) {
	// TODO
	_, _ = exechelper.Do(exechelper.Spec{
		Context: ctx,
		Env:     nil,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	})
}
