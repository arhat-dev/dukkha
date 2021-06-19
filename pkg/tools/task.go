package tools

import (
	"reflect"

	"arhat.dev/dukkha/pkg/field"
)

// TaskType for interface type registration
var TaskType = reflect.TypeOf((*Task)(nil)).Elem()

type Task interface {
	field.Interface

	// Kind of the task, prefixed with tool kind, e.g. golang:build
	Kind() string
}

type BaseTask struct {
	field.BaseField

	Name   string        `yaml:"name"`
	Matrix *MatrixConfig `yaml:"matrix"`
}
