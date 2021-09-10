package plugin

import (
	"encoding/hex"
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strings"

	"arhat.dev/pkg/log"
	"arhat.dev/pkg/randhelper"
	"arhat.dev/rs"
	"github.com/huandu/xstrings"
	"github.com/traefik/yaegi/interp"
	"golang.org/x/mod/module"

	"arhat.dev/dukkha/pkg/dukkha"
)

type SrcRef struct {
	rs.BaseField `yaml:"-"`

	// Source of the single file plugin go code
	Source string `yaml:"source"`

	// Package name and version of the plugin (go module package)
	// (e.g. example.com/foo@v0.1.1)
	// MUST have its dependencies vendored
	Package string `yaml:"package"`
}

func (s *SrcRef) Fetch(spec Spec, cacheDir string) error {
	switch {
	case len(s.Source) != 0:
		// do nothing
		return nil
	case len(s.Package) != 0:
		mod := s.Package
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
	case len(s.Source) != 0:
		src := s.Source
		pkg, err := getPackageOfSource(src)
		if err != nil {
			return fmt.Errorf("invalid source: %w", err)
		}

		if pkg == "main" {
			return fmt.Errorf("invalid source with main package")
		}

		_, err = interp.Eval(src)
		if err != nil {
			return fmt.Errorf("failed to evaluate source code: %w", err)
		}

		return register(spec, pkg, interp)
	case len(s.Package) != 0:
		srcPath := spec.SrcPath(cacheDir)
		pkg, err := getPackageOfDir(srcPath)
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

// TODO: use real factory func type
// 		 currently yaegi doesn't support interface type assignment
// 		 	(see: https://github.com/traefik/yaegi/issues/950)
//
// 			error example cannot convert expression of type *plugin.fooRenderer to type
// 				arhat.dev/dukkha/pkg/dukkha.Renderer

func register(spec Spec, pkg string, pluginPkg *interp.Interpreter) error {
	pkgPrefix := pkg
	if len(pkgPrefix) != 0 {
		pkgPrefix += "."
	}

	switch t := spec.(type) {
	case *RendererSpec:
		funcName := pkgPrefix + getRendererFactoryFuncName(t.DefaultName)

		dukkha.RegisterRenderer(t.DefaultName, func(name string) dukkha.Renderer {
			return NewReflectRenderer(pluginPkg, funcName, name)
		})

		return nil
	case *ToolSpec:
		toolFuncName := pkgPrefix + getToolFactoryFuncName(t.ToolKind)

		tskFuncs := make(map[string]dukkha.TaskCreateFunc)
		for _, tsk := range t.Tasks {
			tskFuncName := pkgPrefix + getTaskFactoryFuncName(t.ToolKind, tsk)
			tskFuncs[tsk] = func(toolName string) dukkha.Task {
				return NewReflectTask(pluginPkg, tskFuncName, toolName)
			}
		}

		dukkha.RegisterTool(dukkha.ToolKind(t.ToolKind), func() dukkha.Tool {
			return NewReflectTool(pluginPkg, toolFuncName)
		})

		for tsk := range tskFuncs {
			dukkha.RegisterTask(
				dukkha.ToolKind(t.ToolKind),
				dukkha.TaskKind(tsk),
				tskFuncs[tsk],
			)
		}
		return nil
	case *TaskSpec:
		funcName := pkgPrefix + getTaskFactoryFuncName(t.ToolKind, t.TaskKind)
		dukkha.RegisterTask(
			dukkha.ToolKind(t.ToolKind),
			dukkha.TaskKind(t.TaskKind),
			func(toolName string) dukkha.Task {
				return NewReflectTask(pluginPkg, funcName, toolName)
			},
		)
		return nil
	default:
		return fmt.Errorf("unknown plugin spec: %T", spec)
	}
}

func newRandomHexBytes() string {
	randBytes, err := randhelper.Bytes(make([]byte, 32))
	if err != nil {
		panic(fmt.Errorf("failed to generate random bytes: %w", err))
	}

	return hex.EncodeToString(randBytes)
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

// getPackageOfSource return go package name of the source code
func getPackageOfSource(src string) (string, error) {
	f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.AllErrors)
	if err != nil {
		return "", err
	}

	if f.Name == nil {
		return "", fmt.Errorf("invalid empty package name")
	}

	return f.Name.Name, nil
}

// getPackageOfDir return go package name of that dir
// if no go source code, return empty string
// if has multiple packages in that dir, their names are compared with `>`,
// 		and that greatest name will be returned
// if has main package in that directory, return `main` only
//
// nolint:deadcode,unused
func getPackageOfDir(dir string) (string, error) {
	pkgs, err := parser.ParseDir(token.NewFileSet(), dir, nil, parser.AllErrors)
	if err != nil {
		return "", err
	}

	ret := ""
	for pkg := range pkgs {
		if strings.HasSuffix(pkg, "_test") {
			continue
		}

		if pkg == "main" {
			return "main", nil
		}

		if pkg > ret {
			ret = pkg
		}
	}

	return ret, nil
}

func evalObjectMethods(
	interp *interp.Interpreter,
	funcCall string,
	requiredMethods []string,
) (map[string]reflect.Value, error) {
	suffix := newRandomHexBytes()
	// tempPkgName := "pkg_" + suffix
	tempPkgName := "plugin"
	varName := "VAR_" + suffix
	_, err := interp.Eval(
		fmt.Sprintf(`
package %s

var %s = %s
`,
			tempPkgName,
			varName, funcCall,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary variable: %w", err)
	}

	v, err := interp.Eval(fmt.Sprintf("%s.%s", tempPkgName, varName))
	if err != nil {
		return nil, fmt.Errorf("failed to reference temporary variable: %w", err)
	}

	rs.InitRecursively(v, dukkha.GlobalInterfaceTypeHandler)

	methods := make(map[string]reflect.Value)
	for _, methodName := range requiredMethods {
		methodFunc, err := interp.Eval(
			fmt.Sprintf(
				`%s.%s.%s`, tempPkgName, varName, methodName,
			),
		)
		if err != nil {
			panic(err)
		}

		methods[methodName] = methodFunc
	}

	return methods, nil
}
