package templateutils

import (
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
)

// Dukkha runtime specific template funcs

func createDukkhaNS(rc dukkha.RenderingContext) *dukkhaNS {
	return &dukkhaNS{rc: rc}
}

type dukkhaNS struct {
	rc dukkha.RenderingContext
}

func (ns *dukkhaNS) SetValue(key string, v interface{}) (interface{}, error) {
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
			"failed to generate new values for key %q: %w",
			key, err,
		)
	}

	err = ns.rc.AddValues(newValues)
	if err != nil {
		return v, fmt.Errorf("failed to add new value: %w", err)
	}

	return v, nil
}
