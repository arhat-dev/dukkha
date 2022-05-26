package templateutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	dukkha_internal "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/pkg/textquery"
	"arhat.dev/rs"
	"github.com/itchyny/gojq"
	"gopkg.in/yaml.v3"
)

// Dukkha runtime specific template funcs

func createDukkhaNS(rc dukkha.RenderingContext) dukkhaNS { return dukkhaNS{rc: rc} }

type dukkhaNS struct{ rc dukkha.RenderingContext }

type cmdOutput = struct {
	Stdout string
	Stderr string
}

// Self runs dukkha command in current process
func (ns dukkhaNS) Self(args ...String) (ret cmdOutput, err error) {
	var stdout, stderr strings.Builder

	flags, err := toStrings(args)
	if err != nil {
		return
	}

	err = dukkha_internal.RunSelf(
		ns.rc.(dukkha.Context),
		ns.rc.Stdin(),
		&stdout,
		&stderr,
		flags...,
	)

	ret.Stdout = stdout.String()
	ret.Stderr = stderr.String()
	return
}

// CrossPlatform checks if doing cross platform job by comparing
// arg[0]/matrix.kernel with arg[1]/host.kernel
// arg[2]/matrix.arch with arg[3]/host.arch
func (ns dukkhaNS) CrossPlatform(args ...String) bool {
	var (
		hostKernel, hostArch     string
		targetKernel, targetArch string
	)

	switch len(args) {
	case 0:
		targetKernel, targetArch = ns.rc.MatrixKernel(), ns.rc.MatrixArch()
		hostKernel, hostArch = ns.rc.HostKernel(), ns.rc.HostArch()
	case 1:
		targetKernel, targetArch = must(toString(args[0])), ns.rc.MatrixArch()
		hostKernel, hostArch = ns.rc.HostKernel(), ns.rc.HostArch()
	case 2:
		targetKernel, targetArch = must(toString(args[0])), ns.rc.MatrixArch()
		hostKernel, hostArch = must(toString(args[1])), ns.rc.HostArch()
	case 3:
		targetKernel, targetArch = must(toString(args[0])), must(toString(args[2]))
		hostKernel, hostArch = must(toString(args[1])), ns.rc.HostArch()
	default:
		targetKernel, targetArch = must(toString(args[0])), must(toString(args[2]))
		hostKernel, hostArch = must(toString(args[1])), must(toString(args[3]))
	}

	return constant.CrossPlatform(targetKernel, targetArch, hostKernel, hostArch)
}

// CacheDir gets DUKKHA_CACHE_DIR
func (ns dukkhaNS) CacheDir() string { return ns.rc.CacheDir() }

// WorkDir gets DUKKHA_WORKDIR
func (ns dukkhaNS) WorkDir() string { return ns.rc.WorkDir() }

// Set is an alias of SetValue
func (ns dukkhaNS) Set(key String, v any) (any, error) { return ns.SetValue(key, v) }

// SetValue set global value
func (ns dukkhaNS) SetValue(key String, v any) (_ any, err error) {
	strKey, err := toString(key)
	if err != nil {
		return
	}

	// parse yaml/json doc when v is string or bytes
	switch t := v.(type) {
	case string:
		v, err = fromYaml(ns.rc, t)
	case []byte:
		v, err = fromYaml(ns.rc, string(t))
	default:
		// do nothing
	}

	if err != nil {
		return v, err
	}

	// TODO: support jq path reference so we can operate on array
	//       entries

	// const newValueJQVarName = "$dukkha_new_value_for_jq"
	// query, err := gojq.Parse(fmt.Sprintf(".%s = %s", key, newValueJQVarName))
	// if err != nil {
	// 	return v, err
	// }
	// _, _, err = textquery.RunQuery(query, newValues, map[string]any{
	// 	newValueJQVarName: v,
	// })
	// if err != nil {
	// 	return v, err
	// }

	newValues := make(map[string]any)

	err = genNewVal(strKey, v, &newValues)
	if err != nil {
		return v, fmt.Errorf(
			"generate new values for key %q: %w",
			key, err,
		)
	}

	err = ns.rc.AddValues(newValues)
	if err != nil {
		return v, fmt.Errorf("bad new value: %w", err)
	}

	return v, nil
}

func genNewVal(key string, value any, ret *map[string]any) error {
	var (
		thisKey string
		nextKey string
	)

	if strings.HasPrefix(key, `"`) {
		key = key[1:]
		quoteIdx := strings.IndexByte(key, '"')
		if quoteIdx < 0 {
			return fmt.Errorf("invalid unclosed quote in string `%s'", key)
		}

		thisKey = key[:quoteIdx]
		nextKey = key[quoteIdx+1:]

		if len(nextKey) == 0 {
			// no more nested maps
			(*ret)[thisKey] = value
			return nil
		}
	} else {
		dotIdx := strings.IndexByte(key, '.')
		if dotIdx < 0 {
			// no more dots, no more nested maps
			(*ret)[key] = value
			return nil
		}

		thisKey = key[:dotIdx]
		nextKey = key[dotIdx+1:]
	}

	newValParent := make(map[string]any)
	(*ret)[thisKey] = newValParent

	return genNewVal(nextKey, value, &newValParent)
}

// JQObj is an alias of YQObj (as YAML is a superset of JSON)
func (ns dukkhaNS) JQObj(args ...any) (_ any, err error) {
	return ns.YQObj(args...)
}

// YQObj is like YQ, but returns object instead of marshaled text
func (ns dukkhaNS) YQObj(args ...any) (_ any, err error) {
	var candidates []any
	err = handleTextQuery(args, ns.rc, func(data any, result []any, queryErr error) error {
		candidates = append(candidates, result...)
		return queryErr
	})
	if err != nil {
		return
	}

	switch len(candidates) {
	case 0:
		return nil, nil
	case 1:
		return candidates[0], nil
	default:
		return candidates, nil
	}
}

// JQ is jq on json string/object, return json string
func (ns dukkhaNS) JQ(args ...any) (_ string, err error) {
	var sb strings.Builder

	err = handleTextQuery(args, ns.rc,
		textquery.CreateResultToTextHandleFuncForJsonOrYaml(&sb, json.Marshal),
	)
	if err != nil {
		return
	}

	return sb.String(), nil
}

// YQ is jq on yaml string/object with multi-doc support, return yaml text
func (ns dukkhaNS) YQ(args ...any) (_ string, err error) {
	var sb strings.Builder

	err = handleTextQuery(args, ns.rc,
		textquery.CreateResultToTextHandleFuncForJsonOrYaml(&sb, yaml.Marshal),
	)
	if err != nil {
		return
	}

	return sb.String(), nil
}

func handleTextQuery(args []any, rc dukkha.RenderingContext, handle textquery.QueryResultHandleFunc) (err error) {
	var (
		query     String
		data      Bytes
		variables map[string]any
	)

	switch n := len(args); n {
	case 0, 1:
		err = fmt.Errorf("at least 2 args expected, got %d", n)
		return
	case 2:
		query, data = args[0], args[1]
	default:
		query, data = args[0], args[n-1]
		switch t := args[1].(type) {
		case map[string]string:
			variables = make(map[string]any, len(t))
			for k, v := range t {
				variables[k] = v
			}
		case map[string]any:
			variables = t
		default:
			err = fmt.Errorf("unsupported variables container type: %T", t)
			return
		}
	}

	q, err := toString(query)
	if err != nil {
		return
	}

	return JQ(rc, q, variables, data, handle)
}

// FromYaml unmarshals single yaml doc into any/map[string]any
func (ns dukkhaNS) FromYaml(v Bytes) (_ any, err error) {
	return fromYaml(ns.rc, v)
}

func fromYaml(rc rs.RenderingHandler, v Bytes) (_ any, err error) {
	inData, inReader, isReader, err := toBytesOrReader(v)
	if err != nil {
		return
	}

	var decode func(out any) error
	if isReader {
		dec := yaml.NewDecoder(inReader)
		decode = dec.Decode
	} else {
		decode = func(out any) error {
			return yaml.Unmarshal(inData, out)
		}
	}

	out := rs.Init(&rs.AnyObject{}, nil).(*rs.AnyObject)
	err = decode(out)
	if err != nil {
		return nil, fmt.Errorf("fromYaml: unmarshal yaml data: %w", err)
	}

	err = out.ResolveFields(rc, -1)
	if err != nil {
		return nil, fmt.Errorf("fromYaml: resolving yaml data: %w", err)
	}

	return out.NormalizedValue(), nil
}

func JQ(
	rc dukkha.RenderingContext,
	query string,
	variables map[string]any,
	v Bytes,
	handleResult textquery.QueryResultHandleFunc,
) error {
	var (
		varNames   []string
		varValues  []any
		errDocIter error

		docIter textquery.DocIterFunc
	)

	inData, inReader, isReader, err := toBytesOrReader(v)
	if err != nil {
		done := false
		docIter = func() (any, bool) {
			if done {
				return nil, false
			}

			done = true
			return v, true
		}
	} else {
		var rd bytes.Reader

		if !isReader {
			rd.Reset(inData)

			inReader = &rd
		}

		dec := yaml.NewDecoder(inReader)

		docIter = func() (any, bool) {
			var obj rs.AnyObject

			_ = rs.Init(&obj, nil)

			errDocIter = dec.Decode(&obj)
			if errDocIter != nil {
				// TODO: return plain text on unexpected error?
				return nil, false
			}

			errDocIter = obj.ResolveFields(rc, -1)
			if errDocIter != nil {
				return nil, false
			}

			return obj.NormalizedValue(), true
		}
	}

	for k, v := range variables {
		varNames = append(varNames, k)
		varValues = append(varValues, v)
	}

	nOptions := 1
	options := [2]gojq.CompilerOption{
		0: gojq.WithEnvironLoader(func() (ret []string) {
			allEnv := rc.Env()
			for k, v := range allEnv {
				ret = append(ret, k+"="+v.Get())
			}

			return
		}),
	}

	if len(varNames) != 0 {
		options[1] = gojq.WithVariables(varNames)
		nOptions = 2
	}

	err = textquery.Query(
		query,
		varValues,
		docIter,
		handleResult,
		options[:nOptions]...,
	)

	if err != nil {
		return err
	}

	return errDocIter
}
