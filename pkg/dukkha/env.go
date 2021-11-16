package dukkha

import (
	"fmt"
	"reflect"

	"arhat.dev/rs"
)

// EnvEntry is a single name/value pair
type EnvEntry struct {
	rs.BaseField `yaml:"-" json:"-"`

	// Name of the entry (in other words, key)
	Name string `yaml:"name"`

	// Value associated to the name
	Value string `yaml:"value"`
}

// Env is a list of name/value pairs (ordered)
type Env []*EnvEntry

// Clone makes a copy of all env values without doing BaseField initialization
func (orig Env) Clone() Env {
	ret := make(Env, 0, len(orig))
	for _, entry := range orig {
		ret = append(ret, &EnvEntry{
			Name:  entry.Name,
			Value: entry.Value,
		})
	}

	return ret
}

// ResolveEnv resolve struct field with `Env` type in struct parent and
// add these env entries to ctx during resolving, so later entries can rely
// on former entries' value
//
// NOTE: parent should be initialized with rs.Init before calling this function
func ResolveEnv(ctx RenderingContext, parent rs.Field, envFieldName, envTagName string) error {
	err := parent.ResolveFields(ctx, 1, envTagName)
	if err != nil {
		return fmt.Errorf("failed to get env overview: %w", err)
	}

	fv := reflect.ValueOf(parent)
	for fv.Kind() == reflect.Ptr {
		fv = fv.Elem()
	}

	// avoid panic
	if fv.Kind() != reflect.Struct {
		return fmt.Errorf("unexpected non struct target: %T", parent)
	}

	env := fv.FieldByName(envFieldName).Interface().(Env)
	for i := range env {
		err = env[i].ResolveFields(ctx, -1)
		if err != nil {
			return fmt.Errorf("failed to resolve env %q: %w", env[i].Name, err)
		}

		ctx.AddEnv(true, env[i])
	}

	return nil
}
