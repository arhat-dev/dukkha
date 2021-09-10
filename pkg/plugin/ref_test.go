package plugin

import (
	"context"
	"reflect"
	"testing"

	"arhat.dev/pkg/log"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"

	_ "embed"
)

var (
	_ dukkha.Renderer = (*fooRenderer)(nil)
	_ dukkha.Tool     = (*toolFoo)(nil)
	_ dukkha.Task     = (*taskFoo)(nil)
)

func init() {
	log.SetDefaultLogger(log.ConfigSet{
		{
			Level:  "verbose",
			Format: "console",
			Destination: log.Destination{
				File: "stderr",
			},
		},
	})
}

// nolint:deadcode,unused,varcheck
//go:embed plugin_example_test.go
var exampleSourceFile string

func TestGetPackageOfSource(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		expectErr bool
		expected  string
	}{
		{
			name:      "Valid Single Package",
			src:       "\n\npackage foo",
			expectErr: false,
			expected:  "foo",
		},
		{
			name:      "Invalid Multi Package",
			src:       "\n\npackage def\n\n package main",
			expectErr: true,
		},
		{
			name:      "Invalid Empty Package",
			src:       "package   ",
			expectErr: true,
		},
		{
			name:      "Invalid No Package",
			src:       "const test = 1",
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pkg, err := getPackageOfSource(test.src)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, pkg)
			}
		})
	}
}

func TestGetPackageOfDir(t *testing.T) {
	tests := []struct {
		name string

		path      string
		expectErr bool
		expected  string
	}{
		{
			name:      "Valid Single Package",
			path:      "testdata/test-module/valid-single-package",
			expectErr: false,
			expected:  "valid_single_package",
		},
		{
			name:      "Valid Multi Package",
			path:      "testdata/test-module/valid-multi-package",
			expectErr: false,
			expected:  "z",
		},
		{
			name:      "Valid No Package",
			path:      "testdata/test-module/valid-no-package",
			expectErr: false,
			expected:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pkg, err := getPackageOfDir(test.path)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, pkg)
			}
		})
	}
}

func TestSrcRef_Register(t *testing.T) {
	tests := []struct {
		name string
		spec Spec
	}{
		{
			name: "Empty Renderer",
			spec: &RendererSpec{
				DefaultName: "foo",
				SrcRef: SrcRef{
					Source: exampleSourceFile,
				},
			},
		},
		{
			name: "Empty Tool",
			spec: &ToolSpec{
				ToolKind: "tool-foo",
				Tasks:    []string{"task-foo", "task-bar"},
				SrcRef: SrcRef{
					Source: exampleSourceFile,
				},
			},
		},
		{
			name: "Empty Task",
			spec: &TaskSpec{
				ToolKind: "tool-bar",
				TaskKind: "task-bar",
				SrcRef: SrcRef{
					Source: exampleSourceFile,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch spec := test.spec.(type) {
			case *RendererSpec:
				err := spec.SrcRef.Register(spec, t.TempDir())
				assert.NoError(t, err)

				rendererIface, err := dukkha.GlobalInterfaceTypeHandler.Create(
					reflect.TypeOf((*dukkha.Renderer)(nil)).Elem(),
					"foo",
				)
				assert.NoError(t, err)

				ret, err := rendererIface.(dukkha.Renderer).RenderYaml(dukkha_test.NewTestContext(context.TODO()), nil)
				assert.NoError(t, err, err)
				assert.EqualValues(t, "HELLO foo", string(ret))
			case *ToolSpec:
				err := spec.SrcRef.Register(spec, t.TempDir())
				assert.NoError(t, err)
			case *TaskSpec:
				err := spec.SrcRef.Register(spec, t.TempDir())
				assert.NoError(t, err)
			default:
				t.Log("unknown spec type", spec)
			}
		})
	}
}
