package templateutils

import (
	"fmt"
	"strings"

	"arhat.dev/pkg/textquery"
	"arhat.dev/rs"
	"github.com/itchyny/gojq"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/third_party/gomplate/conv"
)

func jqObject(query, in interface{}) (interface{}, error) {
	q, err := gojq.Parse(conv.ToString(query))
	if err != nil {
		return nil, err
	}

	ret, _, err := textquery.RunQuery(q, in, nil)
	switch len(ret) {
	case 0:
		return nil, err
	case 1:
		return ret[0], err
	default:
		return ret, err
	}
}

func fromYaml(rc rs.RenderingHandler, v string) (interface{}, error) {
	out := rs.Init(&rs.AnyObject{}, nil).(*rs.AnyObject)
	err := yaml.Unmarshal([]byte(v), out)
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
