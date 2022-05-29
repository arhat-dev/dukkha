package templateutils

import (
	"fmt"
	"reflect"
	"sync"

	"arhat.dev/dukkha/pkg/dukkha"
)

type FuncID uint32

type (
	StaticFuncs      [FuncID_LAST_Static_FUNC]reflect.Value
	ContextualFuncs  [FuncID_LAST_Contextual_FUNC - FuncID_LAST_Static_FUNC]reflect.Value
	PlaceholderFuncs [FuncID_LAST_Placeholder_FUNC - FuncID_LAST_Contextual_FUNC]reflect.Value
)

var (
	execStaticFuncs     StaticFuncs
	execStaticFuncsOnce sync.Once
)

type TemplateFuncFactory = func(rc dukkha.RenderingContext) any

var (
	toolSpecificFuncs = make(map[string]TemplateFuncFactory)
)

func RegisterTemplateFuncs(fm map[string]TemplateFuncFactory) {
	for k, f := range fm {
		if _, ok := toolSpecificFuncs[k]; ok {
			panic(fmt.Sprintf("func %q already registered", k))
		}

		toolSpecificFuncs[k] = f
	}
}

func CreateTemplateFuncs(rc dukkha.RenderingContext) (ret TemplateFuncs) {
	execStaticFuncsOnce.Do(func() {
		for i := range staticFuncs {
			execStaticFuncs[i] = reflect.ValueOf(staticFuncs[i])
		}
	})

	ret.StaticFuncs = &execStaticFuncs
	ret.ContextualFuncs = createContextualFuncs(rc)
	ret.PlaceholderFuncs = new(PlaceholderFuncs)

	if len(toolSpecificFuncs) != 0 {
		ret.other = make(map[string]reflect.Value)
		for k, createTemplateFunc := range toolSpecificFuncs {
			ret.other[k] = reflect.ValueOf(createTemplateFunc(rc))
		}
	}

	return
}

type TemplateFuncs struct {
	*StaticFuncs
	*ContextualFuncs
	*PlaceholderFuncs

	other map[string]reflect.Value
}

var (
	ns_math     mathNS
	ns_archconv archconvNS
	ns_strings  stringsNS
	ns_type     typeNS
	ns_coll     collNS
	ns_cred     credentialNS
	ns_dns      dnsNS
	ns_enc      encNS
	ns_hash     hashNS
	ns_path     pathNS
	ns_re       regexpNS
	ns_sockaddr sockaddrNS
	ns_time     timeNS
	ns_uuid     uuidNS
	ns_golang   golangNS
)

func get_ns_math() mathNS         { return ns_math }
func get_ns_archconv() archconvNS { return ns_archconv }
func get_ns_strings() stringsNS   { return ns_strings }
func get_ns_type() typeNS         { return ns_type }
func get_ns_coll() collNS         { return ns_coll }
func get_ns_cred() credentialNS   { return ns_cred }
func get_ns_dns() dnsNS           { return ns_dns }
func get_ns_enc() encNS           { return ns_enc }
func get_ns_hash() hashNS         { return ns_hash }
func get_ns_path() pathNS         { return ns_path }
func get_ns_re() regexpNS         { return ns_re }
func get_ns_sockaddr() sockaddrNS { return ns_sockaddr }
func get_ns_time() timeNS         { return ns_time }
func get_ns_uuid() uuidNS         { return ns_uuid }

func (f *TemplateFuncs) Has(name string) bool {
	if FuncNameToFuncID(name) != _unknown_template_func {
		return true
	}

	if f.other == nil {
		return false
	}

	_, ok := f.other[name]
	return ok
}

func (f *TemplateFuncs) Override(fid FuncID, newFn reflect.Value) {
	switch {
	case fid > FuncID_LAST_Placeholder_FUNC:
		// TODO: ?
		return
	case fid > FuncID_LAST_Contextual_FUNC:
		// is placeholder func
		f.PlaceholderFuncs[fid-FuncID_LAST_Contextual_FUNC-1] = newFn
	case fid > FuncID_LAST_Static_FUNC:
		// is contextual func
		f.ContextualFuncs[fid-FuncID_LAST_Static_FUNC-1] = newFn
	default:
		// TODO: shall we allow overrding static funcs at all?
		// is static func
		f.StaticFuncs[fid-1] = newFn
	}
}

// GetByID get func value by id, if fid is invalid or unsupported, returned reflect.Value
// will hav Valid() == false
//
// valid fid range is [1, FuncID_LAST_Placeholder_FUNC]
func (f *TemplateFuncs) GetByID(fid FuncID) (ret reflect.Value) {
	if fid == _unknown_template_func {
		return
	}

	switch {
	case fid > FuncID_LAST_Placeholder_FUNC:
		// TODO: ?
		return
	case fid > FuncID_LAST_Contextual_FUNC:
		// is placeholder func
		return f.PlaceholderFuncs[fid-FuncID_LAST_Contextual_FUNC-1]
	case fid > FuncID_LAST_Static_FUNC:
		// is contextual func
		return f.ContextualFuncs[fid-FuncID_LAST_Static_FUNC-1]
	default:
		// is static func
		return f.StaticFuncs[fid-1]
	}
}

// GetByName get func value by name, not limited to predefined funcs
func (f *TemplateFuncs) GetByName(name string) (ret reflect.Value) {
	fid := FuncNameToFuncID(name)
	if fid < FuncID_COUNT {
		return f.GetByID(fid)
	}

	if f.other != nil {
		return f.other[name]
	}

	return
}

func (f *TemplateFuncs) AddOther(funcMap map[string]any) {
	if len(funcMap) == 0 {
		return
	}

	if f.other == nil {
		f.other = make(map[string]reflect.Value)
	}

	for k := range funcMap {
		// overriding predefined funcs is useless
		if FuncNameToFuncID(k) != _unknown_template_func {
			continue
		}

		f.other[k] = reflect.ValueOf(funcMap[k])
	}
}
