package plugin

import (
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strings"

	"arhat.dev/pkg/log"
	"arhat.dev/rs"
	"github.com/huandu/xstrings"
	"github.com/traefik/yaegi/interp"
	"golang.org/x/mod/module"

	"arhat.dev/dukkha/pkg/dukkha"
)

type SrcRef struct {
	rs.BaseField `yaml:"-"`

	// Source of the single file plugin go code
	Source *string `yaml:"source"`

	// Package name and version of the plugin (go module package)
	// (e.g. example.com/foo@v0.1.1)
	// MUST have its dependencies vendored
	Package *string `yaml:"package"`
}

func (s *SrcRef) Fetch(spec Spec, cacheDir string) error {
	switch {
	case s.Source != nil:
		srcContent := *s.Source
		if len(srcContent) == 0 {
			return fmt.Errorf("invalid empty plugin source")
		}

		// do nothing
		return nil
	case s.Package != nil:
		mod := *s.Package
		if len(mod) == 0 {
			return fmt.Errorf("invalid empty plugin git url")
		}

		ver := &module.Version{}
		versionStart := strings.LastIndex(mod, "@")
		if versionStart > 0 {
			ver.Path = mod[:versionStart]
			ver.Version = mod[versionStart+1:]
		} else {
			ver.Path = mod
		}

		pathParts := []string{spec.SrcPath(cacheDir)}
		pathParts = append(pathParts, strings.Split(mod, "/")...)
		dest := filepath.Join(pathParts...)
		_ = dest

		return nil
	default:
		return fmt.Errorf("invalid no plugin fetch method set")
	}
}

func (s *SrcRef) Register(spec Spec, cacheDir string) error {
	logger := log.Log.WithName("plugin").WithFields(log.String("name", spec.Name()))
	interp, err := newInterperter(
		spec.GoPath(cacheDir),
		strings.NewReader(""), // no stdin
		LogWriter(logger, "stdout", "data"),
		LogWriter(logger, "stderr", "data"),
	)
	if err != nil {
		return err
	}

	switch {
	case s.Source != nil:
		pkg := getFileSourcePackage(*s.Source)
		if pkg == "main" {
			return fmt.Errorf("invalid source with main package")
		}

		_, err := interp.Eval(*s.Source)
		if err != nil {
			return fmt.Errorf("failed to evaluate source code: %w", err)
		}

		return register(spec, pkg, interp)
	case s.Package != nil:
		srcPath := spec.SrcPath(cacheDir)
		pkg, err := getDirSourcePackage(srcPath)
		if err != nil {
			return fmt.Errorf("failed to get package source: %w", err)
		}

		if pkg == "main" {
			return fmt.Errorf("invalid source with main package")
		}

		_, err = interp.EvalPath(srcPath)
		if err != nil {
			return fmt.Errorf("failed to evaluate source package: %w", err)
		}

		return register(spec, pkg, interp)
	default:
		return fmt.Errorf("invalid no plugin source")
	}
}

var (
	rendererCreateFuncType = reflect.ValueOf((*dukkha.RendererCreateFunc)(nil)).Elem().Type()
	toolCreateFuncType     = reflect.ValueOf((*dukkha.ToolCreateFunc)(nil)).Elem().Type()
	taskCreateFuncType     = reflect.ValueOf((*dukkha.TaskCreateFunc)(nil)).Elem().Type()
)

func register(spec Spec, pkg string, pluginPkg *interp.Interpreter) error {
	pkgPrefix := pkg
	if len(pkgPrefix) != 0 {
		pkgPrefix += "."
	}

	switch t := spec.(type) {
	case *RendererSpec:
		fFunc, err := pluginPkg.Eval(pkgPrefix + getRendererFactoryFuncName(t.DefaultName))
		if err != nil {
			return fmt.Errorf("failed to find renderer factory func: %w", err)
		}

		if !fFunc.Type().ConvertibleTo(rendererCreateFuncType) {
			return fmt.Errorf(
				"invalid renderer factory func: expect %q, got %q",
				rendererCreateFuncType, fFunc.Type().String(),
			)
		}

		rcf := fFunc.Convert(rendererCreateFuncType).Interface().(dukkha.RendererCreateFunc)
		dukkha.RegisterRenderer(t.DefaultName, rcf)
		return nil
	case *ToolSpec:
		fFunc, err := pluginPkg.Eval(pkgPrefix + getToolFactoryFuncName(t.ToolKind))
		if err != nil {
			return fmt.Errorf("failed to find tool factory func: %w", err)
		}

		if !fFunc.Type().ConvertibleTo(toolCreateFuncType) {
			return fmt.Errorf(
				"invalid tool factory func: expect %q, got %q",
				toolCreateFuncType, fFunc.Type().String(),
			)
		}

		tskFuncs := make(map[string]dukkha.TaskCreateFunc)
		for _, tsk := range t.Tasks {
			tskFunc, err := pluginPkg.Eval(pkgPrefix + getTaskFactoryFuncName(t.ToolKind, tsk))
			if err != nil {
				return fmt.Errorf("failed to find task factory func: %w", err)
			}

			if !tskFunc.Type().ConvertibleTo(taskCreateFuncType) {
				return fmt.Errorf(
					"invalid task factory func: expect %q, got %q",
					taskCreateFuncType, tskFunc.Type().String(),
				)
			}

			tskFuncs[tsk] = tskFunc.Convert(taskCreateFuncType).Interface().(dukkha.TaskCreateFunc)
		}

		rcf := fFunc.Convert(toolCreateFuncType).Interface().(dukkha.ToolCreateFunc)
		dukkha.RegisterTool(dukkha.ToolKind(t.ToolKind), rcf)
		for tsk := range tskFuncs {
			dukkha.RegisterTask(dukkha.ToolKind(t.ToolKind), dukkha.TaskKind(tsk), tskFuncs[tsk])
		}
		return nil
	case *TaskSpec:
		tskFunc, err := pluginPkg.Eval(pkgPrefix + getTaskFactoryFuncName(t.ToolKind, t.TaskKind))
		if err != nil {
			return fmt.Errorf("failed to find task factory func: %w", err)
		}

		if !tskFunc.Type().ConvertibleTo(taskCreateFuncType) {
			return fmt.Errorf(
				"invalid renderer factory func: expect %q, got %q",
				taskCreateFuncType, tskFunc.Type().String(),
			)
		}

		rcf := tskFunc.Convert(taskCreateFuncType).Interface().(dukkha.TaskCreateFunc)
		dukkha.RegisterTask(dukkha.ToolKind(t.ToolKind), dukkha.TaskKind(t.TaskKind), rcf)
		return nil
	default:
		return fmt.Errorf("unknown plugin spec: %T", spec)
	}
}

func getRendererFactoryFuncName(defaultName string) string {
	return getFactoryFuncName("NewRenderer_", defaultName)
}

func getToolFactoryFuncName(toolKind string) string {
	return getFactoryFuncName("NewTool_", toolKind)
}

func getTaskFactoryFuncName(toolKind, taskKind string) string {
	return getFactoryFuncName("NewTask_", toolKind+"_"+taskKind)
}

func getFactoryFuncName(prefix, name string) string {
	return prefix + xstrings.ToSnakeCase(name)
}

func getFileSourcePackage(src string) string {
	f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.PackageClauseOnly)
	if err != nil {
		return ""
	}

	if f.Name == nil {
		return ""
	}

	return f.Name.Name
}

// nolint:deadcode,unused
func getDirSourcePackage(dir string) (string, error) {
	pkgs, err := parser.ParseDir(token.NewFileSet(), dir, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}

	ret := ""
	for pkg := range pkgs {
		if strings.HasSuffix(pkg, "_test") {
			continue
		}

		if pkg == "main" {
			ret = pkg
		} else if ret != "main" {
			ret = pkg
		}
	}

	return ret, nil
}
