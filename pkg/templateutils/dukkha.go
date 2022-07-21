package templateutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"arhat.dev/pkg/textquery"
	"arhat.dev/rs"
	"github.com/itchyny/gojq"
	"gopkg.in/yaml.v3"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

// dukkha runtime specific template funcs

func createDukkhaNS(rc dukkha.RenderingContext) dukkhaNS { return dukkhaNS{rc: rc} }

type dukkhaNS struct{ rc dukkha.RenderingContext }

type cmdOutput = struct {
	Stdout string
	Stderr string
}

// Self runs dukkha command in current process
// TODO: support writer as the second last argument
func (ns dukkhaNS) Self(args ...String) (ret cmdOutput, err error) {
	var stdin io.Reader

	nArgs := len(args)
	if nArgs != 0 {
		var ok bool
		stdin, ok = args[nArgs-1].(io.Reader)
		if ok {
			nArgs--
		} else {
			stdin = ns.rc.Stdin()
		}
	} else {
		stdin = ns.rc.Stdin()
	}

	var stdout, stderr strings.Builder

	flags, err := toStrings(args[:nArgs])
	if err != nil {
		return
	}

	err = di.RunSelf(
		ns.rc.(dukkha.Context),
		stdin,
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
		v, err = ns.FromYaml(t)
	case []byte:
		v, err = ns.FromYaml(t)
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

func newJSONDecoder(r io.Reader) DataDecoder {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec
}

func newYAMLDecoder(r io.Reader) DataDecoder {
	dec := yaml.NewDecoder(r)
	return dec
}

// JQObj is like JQ, but returns object instead of marshaled text
func (ns dukkhaNS) JQObj(args ...any) (_ any, err error) {
	var ret []any
	err = handleTextQuery(ns.rc, args,
		newJSONDecoder,
		func(data any, result []any, queryErr error) error {
			ret = append(ret, result...)
			return queryErr
		},
	)
	if err != nil {
		return
	}

	switch len(ret) {
	case 0:
		return nil, nil
	case 1:
		return ret[0], nil
	default:
		return ret, nil
	}
}

// YQObj is like YQ, but returns object instead of marshaled text
func (ns dukkhaNS) YQObj(args ...any) (_ any, err error) {
	var ret []any
	err = handleTextQuery(ns.rc, args,
		newYAMLDecoder,
		func(data any, result []any, queryErr error) error {
			ret = append(ret, result...)
			return queryErr
		},
	)
	if err != nil {
		return
	}

	switch len(ret) {
	case 0:
		return nil, nil
	case 1:
		return ret[0], nil
	default:
		return ret, nil
	}
}

// JQ is jq on json string/object with json stream support, return json text
func (ns dukkhaNS) JQ(args ...any) (_ string, err error) {
	var sb strings.Builder

	err = handleTextQuery(ns.rc, args,
		newJSONDecoder,
		textquery.CreateResultToTextHandleFuncForJsonOrYaml(&sb, json.Marshal),
	)
	if err != nil {
		sb.Reset()
		return
	}

	return sb.String(), nil
}

// YQ is jq on yaml string/object with multi-doc support, return yaml text
func (ns dukkhaNS) YQ(args ...any) (_ string, err error) {
	var sb strings.Builder

	err = handleTextQuery(ns.rc, args,
		newYAMLDecoder,
		textquery.CreateResultToTextHandleFuncForJsonOrYaml(&sb, yaml.Marshal),
	)
	if err != nil {
		sb.Reset()
		return
	}

	return sb.String(), nil
}

func handleTextQuery(
	rc dukkha.RenderingContext,
	args []any, // TODO: support writer as the second last argument
	newDecoder func(io.Reader) DataDecoder,
	handle textquery.QueryResultHandleFunc,
) (err error) {
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
			err = fmt.Errorf("unsupported type of variables %T: expecting map[stirng]{string, any}", t)
			return
		}
	}

	q, err := toString(query)
	if err != nil {
		return
	}

	return JQ(rc, data, JQOptions{
		Query: q,

		// TODO: support optional resolving
		ResolveRenderingSuffixBeforeQueryStart: false,

		Variables:    variables,
		NewDecoder:   newDecoder,
		HandleResult: handle,
	})
}

// FromYaml unmarshals single yaml doc into []any/map[string]any
func (ns dukkhaNS) FromYaml(v Bytes) (_ any, err error) {
	return FromText(ns.rc, v,
		newYAMLDecoder,
		yaml.Unmarshal,
	)
}

// FromJson
// nolint:revive
func (ns dukkhaNS) FromJson(v Bytes) (any, error) {
	return FromText(ns.rc, v,
		newJSONDecoder,
		json.Unmarshal,
	)
}

// TODO: support optional resolving
func FromText(
	rc dukkha.RenderingContext,
	v Bytes,
	newDecoder func(io.Reader) DataDecoder,
	unmarshal func(data []byte, out any) error,
) (_ any, err error) {
	inData, inReader, isReader, err := toBytesOrReader(v)
	if err != nil {
		return
	}

	var decode func(out any) error
	if isReader {
		dec := newDecoder(inReader)
		decode = dec.Decode
	} else {
		decode = func(out any) error {
			return unmarshal(inData, out)
		}
	}

	var out rs.AnyObject

	_ = rs.InitAny(&out, nil)
	err = decode(&out)
	if err != nil {
		return nil, fmt.Errorf("fromX: unamrshal data: %w", err)
	}

	err = out.ResolveFields(rc, -1)
	if err != nil {
		return nil, fmt.Errorf("fromX: resolve data: %w", err)
	}

	return out.NormalizedValue(), nil
}

type DataDecoder interface {
	Decode(any) error
}

type JQOptions struct {
	// Query is the jq expression
	// REQUIRED
	Query string

	// resolve foo@xxx in decoded data to foo before running jq over it
	ResolveRenderingSuffixBeforeQueryStart bool

	// Variables provided when running jq
	//
	// OPTIONAL
	Variables map[string]any

	// NewDecoder creates a data decoder
	//
	// the created decoder should support both rs.AnyObject and rs.AnyObjectMap (usually {json, yaml}.NewDecoder)
	// when ResolveRenderingSuffixBeforeQueryStart is set to true
	//
	// OPTIONAL, defaults to yaml.NewDecoder
	NewDecoder func(io.Reader) DataDecoder

	// HandleResult called for each unmarshaled object
	//
	// REQUIRED
	HandleResult textquery.QueryResultHandleFunc
}

// JQ evaluates jq expression over v
func JQ(
	rc dukkha.RenderingContext,
	v Bytes,
	opts JQOptions,
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

		var dec DataDecoder

		if opts.NewDecoder == nil {
			dec = yaml.NewDecoder(inReader)
		} else {
			dec = opts.NewDecoder(inReader)
		}

		if opts.ResolveRenderingSuffixBeforeQueryStart {
			docIter = func() (any, bool) {
				var obj rs.AnyObject

				_ = rs.InitAny(&obj, nil)

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
		} else {
			docIter = func() (any, bool) {
				var obj any
				errDocIter = dec.Decode(&obj)
				if errDocIter != nil {
					return nil, false
				}

				return obj, true
			}
		}
	}

	for k, v := range opts.Variables {
		varNames = append(varNames, k)
		varValues = append(varValues, v)
	}

	nOptions := 1
	options := [2]gojq.CompilerOption{
		0: gojq.WithEnvironLoader(func() (ret []string) {
			allEnv := rc.Env()
			for k, v := range allEnv {
				ret = append(ret, k+"="+v.GetLazyValue())
			}

			return
		}),
	}

	if len(varNames) != 0 {
		options[1] = gojq.WithVariables(varNames)
		nOptions = 2
	}

	err = textquery.Query(
		opts.Query,
		varValues,
		docIter,
		opts.HandleResult,
		options[:nOptions]...,
	)

	if err != nil {
		return err
	}

	if errors.Is(errDocIter, io.EOF) {
		errDocIter = nil
	}

	return errDocIter
}
