package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"

	_ "embed"
)

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
					Source: &exampleSourceFile,
				},
			},
		},
		{
			name: "Empty Tool",
			spec: &ToolSpec{
				ToolKind: "foo-tool",
				Tasks:    []string{"foo-task", "bar-task"},
				SrcRef: SrcRef{
					Source: &exampleSourceFile,
				},
			},
		},
		{
			name: "Empty Task",
			spec: &TaskSpec{
				ToolKind: "bar-tool",
				TaskKind: "bar-task",
				SrcRef: SrcRef{
					Source: &exampleSourceFile,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			switch spec := test.spec.(type) {
			case *RendererSpec:
				err = spec.SrcRef.Register(spec, t.TempDir())
			case *ToolSpec:
				err = spec.SrcRef.Register(spec, t.TempDir())
			case *TaskSpec:
				err = spec.SrcRef.Register(spec, t.TempDir())
			default:
				t.Log("unknown spec type", spec)
			}

			assert.NoError(t, err)
		})
	}
}
