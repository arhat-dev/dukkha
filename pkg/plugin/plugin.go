package plugin

import (
	"fmt"
	"io"
	"path/filepath"
	"reflect"

	"arhat.dev/rs"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

var (
	dukkhaSymbols = interp.Exports{
		"arhat.dev/dukkha/pkg/dukkha": map[string]reflect.Value{
			"ToolName":        reflect.ValueOf((*dukkha.ToolName)(nil)).Elem(),
			"TaskName":        reflect.ValueOf((*dukkha.TaskName)(nil)).Elem(),
			"ArbitraryValues": reflect.ValueOf((*dukkha.ArbitraryValues)(nil)).Elem(),

			"Renderer": reflect.ValueOf((*dukkha.Renderer)(nil)).Elem(),
			"Tool":     reflect.ValueOf((*dukkha.Tool)(nil)).Elem(),
			"Task":     reflect.ValueOf((*dukkha.Task)(nil)).Elem(),

			"ToolKey": reflect.ValueOf((*dukkha.ToolKey)(nil)).Elem(),
			"TaskKey": reflect.ValueOf((*dukkha.TaskKey)(nil)).Elem(),

			"RendererCreateFunc": reflect.ValueOf((*dukkha.RendererCreateFunc)(nil)).Elem(),
			"ToolCreateFunc":     reflect.ValueOf((*dukkha.ToolCreateFunc)(nil)).Elem(),
			"TaskCreateFunc":     reflect.ValueOf((*dukkha.TaskCreateFunc)(nil)).Elem(),
		},
		"arhat.dev/dukkha/pkg/tools": map[string]reflect.Value{
			"BaseTask":        reflect.ValueOf((*tools.BaseTask)(nil)).Elem(),
			"BaseTool":        reflect.ValueOf((*tools.BaseTool)(nil)).Elem(),
			"TaskExecRequest": reflect.ValueOf((*tools.TaskExecRequest)(nil)).Elem(),
			"Action":          reflect.ValueOf((*tools.Action)(nil)).Elem(),
			"TaskHooks":       reflect.ValueOf((*tools.TaskHooks)(nil)).Elem(),
		},
		"arhat.dev/dukkha/pkg/sliceutils": map[string]reflect.Value{
			"NewStrings":      reflect.ValueOf(sliceutils.NewStrings),
			"FormatStringMap": reflect.ValueOf(sliceutils.FormatStringMap),
		},
		"arhat.dev/dukkha/pkg/constant": map[string]reflect.Value{
			"GetOciOS":   reflect.ValueOf(constant.GetOciOS),
			"ARCH_AMD64": reflect.ValueOf(constant.ARCH_AMD64),
		},
	}
)

func init() {
	dukkhaSymbols["arhat.dev/dukkha/pkg/plugin"] = map[string]reflect.Value{
		"Symbols": reflect.ValueOf(dukkhaSymbols),
	}
}

func newInterperter(
	goPath string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) (*interp.Interpreter, error) {
	t := interp.New(interp.Options{
		GoPath: goPath,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})

	err := t.Use(stdlib.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use std libraries: %w", err)
	}

	err = t.Use(syscall.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use syscall libraries: %w", err)
	}

	err = t.Use(unsafe.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use unsafe libraries: %w", err)
	}

	err = t.Use(unrestricted.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use unrestricted libraries: %w", err)
	}

	err = t.Use(dukkhaSymbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use dukkha libraries: %w", err)
	}

	return t, nil
}

type Spec interface {
	Name() string
	GoPath(cacheDir string) string
	SrcPath(cacheDir string) string
}

func goPathSrc(goPath string) string {
	return filepath.Join(goPath, "src")
}

// RendererSpec defines a new renderer or orverrides
// existing renderer
type RendererSpec struct {
	rs.BaseField `yaml:"-"`

	DefaultName string `yaml:"name"`

	SrcRef `yaml:",inline"`
}

func (s *RendererSpec) Name() string                   { return "renderer-" + s.DefaultName }
func (s *RendererSpec) GoPath(cacheDir string) string  { return filepath.Join(cacheDir, s.Name()) }
func (s *RendererSpec) SrcPath(cacheDir string) string { return goPathSrc(s.GoPath(cacheDir)) }

// ToolSpec defines a new tool or overrides existing tool
type ToolSpec struct {
	rs.BaseField `yaml:"-"`

	ToolKind string   `yaml:"tool"`
	Tasks    []string `yaml:"tasks"`

	SrcRef `yaml:",inline"`
}

func (s *ToolSpec) Name() string                   { return "tool-" + s.ToolKind }
func (s *ToolSpec) GoPath(cacheDir string) string  { return filepath.Join(cacheDir, s.Name()) }
func (s *ToolSpec) SrcPath(cacheDir string) string { return goPathSrc(s.GoPath(cacheDir)) }

// TaskSpec presents a task definition for existing tool
type TaskSpec struct {
	rs.BaseField `yaml:"-"`

	ToolKind string `yaml:"tool"`
	TaskKind string `yaml:"task"`

	SrcRef `yaml:",inline"`
}

func (s *TaskSpec) Name() string                   { return "task-" + s.ToolKind + "-" + s.TaskKind }
func (s *TaskSpec) GoPath(cacheDir string) string  { return filepath.Join(cacheDir, s.Name()) }
func (s *TaskSpec) SrcPath(cacheDir string) string { return goPathSrc(s.GoPath(cacheDir)) }
