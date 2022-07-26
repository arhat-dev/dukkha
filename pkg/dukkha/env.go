package dukkha

import (
	"fmt"
	"reflect"

	"arhat.dev/rs"
)

// NameValueEntry is a single name/value pair
type NameValueEntry struct {
	rs.BaseField `yaml:"-" json:"-"`

	// Name of the entry (in other words, key)
	Name string `yaml:"name"`

	// Value associated to the name
	Value string `yaml:"value"`
}

// NameValueList is a list of name/value pairs
type NameValueList []*NameValueEntry

// Clone makes a copy of all env values without doing BaseField initialization
func (orig NameValueList) Clone() NameValueList {
	ret := make(NameValueList, 0, len(orig))
	for _, entry := range orig {
		ret = append(ret, &NameValueEntry{
			Name:  entry.Name,
			Value: entry.Value,
		})
	}

	return ret
}

// ResolveAndAddEnv resolves a NameValueList typed field in parent and add these entries
// into ctx as environment variables, later entries can rely on former entries' value
//
// NOTE: parent should be initialized with rs.Init before calling this function
func ResolveAndAddEnv(ctx RenderingContext, parent rs.Field, listFieldName, listTagName string) error {
	err := parent.ResolveFields(ctx, 1, listTagName)
	if err != nil {
		return fmt.Errorf("gain overview of env: %w", err)
	}

	fv := reflect.ValueOf(parent)
	for fv.Kind() == reflect.Ptr {
		fv = fv.Elem()
	}

	// avoid panic
	if fv.Kind() != reflect.Struct {
		return fmt.Errorf("unexpected non struct target: %T", parent)
	}

	env := fv.FieldByName(listFieldName).Interface().(NameValueList)
	for i := range env {
		err = env[i].ResolveFields(ctx, -1)
		if err != nil {
			return fmt.Errorf("resolving env %q: %w", env[i].Name, err)
		}

		ctx.AddEnv(true, env[i])
	}

	return nil
}
