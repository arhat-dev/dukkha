package dukkha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/textquery"
	"arhat.dev/rs"
	"arhat.dev/tlang"
	"github.com/itchyny/gojq"
	"mvdan.cc/sh/v3/expand"
)

type RenderingContext interface {
	context.Context
	expand.Environ
	rs.InterfaceTypeHandler
	rs.RenderingHandler
	EnvValues

	// FS returns the filesystem with cwd set to DUKKHA_WORKDIR
	FS() *fshelper.OSFS

	// AddValues will merge provided values into existing values
	AddValues(values map[string]any) error

	Env() map[string]tlang.LazyValueType[string]

	Values() map[string]any

	GlobalCacheFS(subdir string) *fshelper.OSFS

	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	SetStdIO(in io.Reader, out, err io.Writer)
}

func newContextRendering(
	ctx contextStd,
	ifaceTypeHandler rs.InterfaceTypeHandler,
	globalEnv map[string]tlang.LazyValueType[string],
) contextRendering {
	envValues := newEnvValues(globalEnv)
	return contextRendering{
		contextStd: ctx,

		envValues: envValues,

		ifaceTypeHandler: ifaceTypeHandler,
		renderers:        make(map[string]Renderer),
		values:           make(map[string]any),

		fs: lazyEnsuredSubFS(fshelper.NewOSFS(false, func() (string, error) {
			return envValues.WorkDir(), nil
		}), "."),
		cacheFS: lazyEnsuredSubFS(fshelper.NewOSFS(false, func() (string, error) {
			return envValues.CacheDir(), nil
		}), "."),
	}
}

// contextRendering is the core context of dukkhaContext, handling most of
// rendering related jobs
//
// it MUST be the first element in dukkhaContext, and MUST be derived together with dukkhaContext
type contextRendering struct {
	contextStd
	envValues

	ifaceTypeHandler rs.InterfaceTypeHandler
	renderers        map[string]Renderer

	values map[string]any

	// nolint:revive
	_VALUE any

	fs      *fshelper.OSFS
	cacheFS *fshelper.OSFS

	stdin          io.Reader
	stdout, stderr io.Writer
}

func (c *contextRendering) clone(newCtx contextStd, separateEnv bool) contextRendering {
	vals := c.envValues
	if separateEnv {
		vals = c.envValues.clone()
	}

	return contextRendering{
		contextStd: newCtx,

		envValues: vals,
		renderers: c.renderers,

		// values are global scoped, DO NOT deep copy in any case
		values: c.values,

		fs: lazyEnsuredSubFS(fshelper.NewOSFS(false, func() (string, error) {
			return vals.WorkDir(), nil
		}), "."),
		cacheFS: lazyEnsuredSubFS(fshelper.NewOSFS(false, func() (string, error) {
			return vals.CacheDir(), nil
		}), "."),

		stdin:  c.stdin,
		stdout: c.stdout,
		stderr: c.stderr,
	}
}

func (c *contextRendering) Stdin() io.Reader {
	if c.stdin == nil {
		return os.Stdin
	}

	return c.stdin
}

func (c *contextRendering) Stdout() io.Writer {
	if c.stdout == nil {
		return os.Stdout
	}

	return c.stdout
}

func (c *contextRendering) Stderr() io.Writer {
	if c.stderr == nil {
		return os.Stderr
	}

	return c.stderr
}

func (c *contextRendering) SetStdIO(stdin io.Reader, stdout, stderr io.Writer) {
	c.stdin, c.stdout, c.stderr = stdin, stdout, stderr
}

func (c *contextRendering) FS() *fshelper.OSFS { return c.fs }

func (c *contextRendering) GlobalCacheFS(subdir string) *fshelper.OSFS {
	return lazyEnsuredSubFS(c.cacheFS, subdir)
}

func lazyEnsuredSubFS(ofs *fshelper.OSFS, subdir string) *fshelper.OSFS {
	if path.IsAbs(subdir) || filepath.IsAbs(subdir) {
		panic(fmt.Errorf("expecting relative path, got %q", subdir))
	}

	return fshelper.NewOSFS(false, func() (string, error) {
		pc, _, _, ok := runtime.Caller(2)
		if ok {
			callerFunc := runtime.FuncForPC(pc).Name()
			switch {
			case strings.HasSuffix(callerFunc, "Abs"):
				// do nothing but return base dir when calling Abs
				return ofs.Abs(subdir)
			default:
			}
		}

		_, err := ofs.Stat(subdir)
		if err == nil {
			return ofs.Abs(subdir)
		}

		// we can only handle dir not exist error

		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}

		err = ofs.MkdirAll(subdir, 0755)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			return "", err
		}

		return ofs.Abs(subdir)
	})
}

// Env returns all environment variables available
// global environment variables are always kept
func (c *contextRendering) Env() map[string]tlang.LazyValueType[string] {
	for k, v := range c.envValues.globalEnv {
		c.envValues.env[k] = v
	}

	return c.envValues.env
}

func (c *contextRendering) AddValues(values map[string]any) error {
	mergedValues, err := rs.MergeMap(c.values, values, false, false)
	if err != nil {
		return err
	}

	c.values = mergedValues
	return nil
}

func (c *contextRendering) Values() map[string]any {
	return c.values
}

// RenderYaml implements rs.RenderingHandler
func (c *contextRendering) RenderYaml(renderer string, rawData any) ([]byte, error) {
	var attributes []RendererAttribute
	attrStart := strings.LastIndexByte(renderer, '#')
	if attrStart != -1 {
		for _, attr := range strings.Split(renderer[attrStart+1:], ",") {
			attributes = append(attributes, RendererAttribute(strings.TrimSpace(attr)))
		}

		renderer = renderer[:attrStart]
	}

	v, ok := c.renderers[renderer]
	if !ok {
		return nil, fmt.Errorf("renderer %q not found", renderer)
	}

	return v.RenderYaml((*dukkhaContext)(unsafe.Pointer(c)), rawData, attributes)
}

func (c *contextRendering) Create(typ reflect.Type, yamlKey string) (any, error) {
	return c.ifaceTypeHandler.Create(typ, yamlKey)
}

func (c *contextRendering) AddRenderer(name string, r Renderer) {
	c.renderers[name] = r
}

func (c *contextRendering) AllRenderers() map[string]Renderer {
	return c.renderers
}

// Get implements expand.Environ
func (c *contextRendering) Get(name string) expand.Variable {
	v, exists := c.globalEnv[name]
	if exists {
		return createVariable(v.GetLazyValue())
	}

	v, exists = c.env[name]
	if exists {
		return createVariable(v.GetLazyValue())
	}

	// non existing env

	// TODO: can we remove all these cases? (except the default case)
	switch name {
	case "IFS":
		v = tlang.ImmediateString(" \t\n")
	case "OPTIND":
		v = tlang.ImmediateString("1")
	case "PWD":
		v = tlang.ImmediateString(c.WorkDir())
	case "UID":
		v = tlang.ImmediateString(
			strconv.FormatInt(int64(os.Getuid()), 10),
		)
	case "GID":
		v = tlang.ImmediateString(
			strconv.FormatInt(int64(os.Getegid()), 10),
		)
	default:
		kind := expand.Unset
		if strings.HasPrefix(name, valuesEnvPrefix) {
			valRef := strings.TrimPrefix(name, valuesEnvPrefix)

			query, err := gojq.Parse("." + valRef)
			if err != nil {
				goto ret
			}

			result, err := textquery.RunQuery(query, c.values, nil)
			if err != nil {
				goto ret
			}

			if len(result) == 0 {
				goto ret
			}

			kind = expand.String
			v = tlang.ImmediateString(textquery.MarshalJsonOrYamlQueryResult(result, json.Marshal))
		}

	ret:
		str := ""
		if v != nil {
			str = v.GetLazyValue()
		}
		return expand.Variable{
			Local:    false,
			Exported: true,
			ReadOnly: false,
			Kind:     kind,
			Str:      str,
		}
	}

	return createVariable(v.GetLazyValue())
}

// Each implements expand.Environ
func (c *contextRendering) Each(do func(name string, vr expand.Variable) bool) {
	visited := make(map[string]struct{})

	for k, v := range c.globalEnv {
		visited[k] = struct{}{}

		if !do(k, createVariable(v.GetLazyValue())) {
			return
		}
	}

	for k, v := range c.env {
		// do not override
		if _, ok := visited[k]; ok {
			continue
		}

		if !do(k, createVariable(v.GetLazyValue())) {
			return
		}
	}

	values, _ := genEnvFromValues(c.values)
	for k, v := range values {
		if !do(k, v) {
			return
		}
	}
}

const valuesEnvPrefix = "values."

func genEnvFromValues(values map[string]any) (map[string]expand.Variable, error) {
	out := make(map[string]expand.Variable)
	for k, v := range values {
		err := genEnvFromInterface(valuesEnvPrefix+k, v, &out)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func genEnvFromInterface(prefix string, v any, out *map[string]expand.Variable) error {
	switch t := v.(type) {
	case map[string]any:
		dataBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}

		(*out)[prefix] = createVariable(string(dataBytes))

		for k, v := range t {
			err = genEnvFromInterface(prefix+"."+k, v, out)
			if err != nil {
				return err
			}
		}
	case string:
		(*out)[prefix] = createVariable(t)
	case []byte:
		(*out)[prefix] = createVariable(string(t))
	default:
		dataBytes, err := json.Marshal(t)
		if err != nil {
			return err
		}

		(*out)[prefix] = createVariable(string(dataBytes))
	}

	return nil
}

// createVariable for embedded shell, if exists is false, will lookup values for the name
func createVariable(value string) expand.Variable {
	// TODO: set kind for lists
	return expand.Variable{
		Local:    false,
		Exported: true,
		ReadOnly: false,
		Kind:     expand.String,
		Str:      value,
	}
}
