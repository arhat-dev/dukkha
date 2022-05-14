package templateutils

import (
	"fmt"
	"strings"

	"arhat.dev/rs"
	"gopkg.in/yaml.v3"
)

func fromYaml(rc rs.RenderingHandler, v Bytes) (any, error) {
	out := rs.Init(&rs.AnyObject{}, nil).(*rs.AnyObject)
	err := yaml.Unmarshal(toBytes(v), out)
	if err != nil {
		return nil, fmt.Errorf("unmarshal yaml data\n\n%s\n\nerr: %w", v, err)
	}

	err = out.ResolveFields(rc, -1)
	if err != nil {
		return nil, fmt.Errorf(
			"resolving yaml data\n\n%s\n\nerr: %w",
			v, err,
		)
	}

	return out.NormalizedValue(), nil
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
