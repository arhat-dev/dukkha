package templateutils

import (
	"fmt"
	"strings"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"
)

func fromYaml(rc rs.RenderingHandler, v string) (interface{}, error) {
	out := new(rs.AnyObject)
	err := yaml.Unmarshal([]byte(v), out)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml data\n\n%s\n\nerr: %w", v, err)
	}

	switch out.Value().(type) {
	case string, []byte:
		return v, nil
	}

	err = out.ResolveFields(rc, -1)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve yaml data\n\n%s\n\nerr: %w",
			v, err,
		)
	}

	rawBytes, err := yamlhelper.ToYamlBytes(out)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to marshal back resolved value of \n\n%s\n\nerr: %w",
			v, err,
		)
	}

	var ret interface{}
	err = yaml.Unmarshal(rawBytes, &ret)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal resolved value of \n\n%s\n\n as plain object: %w",
			v, err,
		)
	}

	return ret, nil
}

func genNewVal(key string, value interface{}, ret *map[string]interface{}) error {
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

	newValParent := make(map[string]interface{})
	(*ret)[thisKey] = newValParent

	return genNewVal(nextKey, value, &newValParent)
}
