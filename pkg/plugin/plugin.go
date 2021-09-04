package plugin

import (
	"fmt"
	"io"
	"path/filepath"

	"arhat.dev/rs"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"

	constant_symbols "arhat.dev/dukkha/pkg/constant/symbols"
	dukkha_symbols "arhat.dev/dukkha/pkg/dukkha/symbols"
	matrix_symbols "arhat.dev/dukkha/pkg/matrix/symbols"
	renderer_symbols "arhat.dev/dukkha/pkg/renderer/symbols"
	sliceutils_symbols "arhat.dev/dukkha/pkg/sliceutils/symbols"
	templateutils_symbols "arhat.dev/dukkha/pkg/templateutils/symbols"
	tools_symbols "arhat.dev/dukkha/pkg/tools/symbols"
	rs_symbols "arhat.dev/dukkha/third_party/rs/symbols"
	yaml_symbols "arhat.dev/dukkha/third_party/yaml.v3/symbols"
)

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

	err = t.Use(dukkha_symbols.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use dukkha symbols: %w", err)
	}

	err = t.Use(matrix_symbols.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use matrix symbols: %w", err)
	}

	err = t.Use(renderer_symbols.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use renderer symbols: %w", err)
	}

	err = t.Use(sliceutils_symbols.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use sliceutils symbols: %w", err)
	}

	err = t.Use(templateutils_symbols.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use templateutils symbols: %w", err)
	}

	err = t.Use(tools_symbols.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use tools symbols: %w", err)
	}

	err = t.Use(constant_symbols.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use constant symbols: %w", err)
	}

	err = t.Use(yaml_symbols.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use yaml symbols: %w", err)
	}

	err = t.Use(rs_symbols.Symbols)
	if err != nil {
		return nil, fmt.Errorf("unable to use rs symbols: %w", err)
	}

	t.ImportUsed()

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
