package templateutils

import (
	"fmt"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

// Dukkha runtime specific template funcs

func createDukkhaNS(rc dukkha.RenderingContext) dukkhaNS { return dukkhaNS{rc: rc} }

type dukkhaNS struct{ rc dukkha.RenderingContext }

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
		targetKernel, targetArch = toString(args[0]), ns.rc.MatrixArch()
		hostKernel, hostArch = ns.rc.HostKernel(), ns.rc.HostArch()
	case 2:
		targetKernel, targetArch = toString(args[0]), ns.rc.MatrixArch()
		hostKernel, hostArch = toString(args[1]), ns.rc.HostArch()
	case 3:
		targetKernel, targetArch = toString(args[0]), toString(args[2])
		hostKernel, hostArch = toString(args[1]), ns.rc.HostArch()
	default:
		targetKernel, targetArch = toString(args[0]), toString(args[2])
		hostKernel, hostArch = toString(args[1]), toString(args[3])
	}

	return constant.CrossPlatform(targetKernel, targetArch, hostKernel, hostArch)
}

// CacheDir gets DUKKHA_CACHE_DIR
func (ns dukkhaNS) CacheDir() string { return ns.rc.CacheDir() }

// WorkDir gets DUKKHA_WORKDIR
func (ns dukkhaNS) WorkDir() string { return ns.rc.WorkDir() }

// Set is an alias of SetValue
func (ns dukkhaNS) Set(key string, v interface{}) (interface{}, error) { return ns.SetValue(key, v) }

// SetValue set global value
func (ns dukkhaNS) SetValue(key string, v interface{}) (interface{}, error) {
	var err error
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
	// _, _, err = textquery.RunQuery(query, newValues, map[string]interface{}{
	// 	newValueJQVarName: v,
	// })
	// if err != nil {
	// 	return v, err
	// }

	newValues := make(map[string]interface{})

	err = genNewVal(key, v, &newValues)
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
