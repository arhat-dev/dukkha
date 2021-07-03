package textquery

import (
	"fmt"
	"strconv"

	"github.com/itchyny/gojq"
)

func Query(
	query string,
	input []byte,
	unmarshalFunc func([]byte, interface{}) error,
	marshalFunc func(in interface{}) ([]byte, error),
) (string, error) {
	q, err := gojq.Parse(query)
	if err != nil {
		return "", fmt.Errorf("failed to parse query: %w", err)
	}

	strData, err := strconv.Unquote(string(input))
	if err == nil {
		input = []byte(strData)
	}

	var data interface{}
	err = unmarshalFunc(input, &data)
	if err != nil {
		// plain text data
		data = input
	}

	result, _, err := RunQuery(q, data, nil)
	return HandleQueryResult(result, marshalFunc), err
}

func RunQuery(
	query *gojq.Query,
	data interface{},
	kvPairs map[string]interface{},
) ([]interface{}, bool, error) {
	var iter gojq.Iter

	if len(kvPairs) == 0 {
		iter = query.Run(data)
	} else {
		var (
			keys   []string
			values []interface{}
		)
		for k, v := range kvPairs {
			keys = append(keys, k)
			values = append(values, v)
		}

		code, err := gojq.Compile(query, gojq.WithVariables(keys))
		if err != nil {
			return nil, false, fmt.Errorf("failed to compile query with variables: %w", err)
		}

		iter = code.Run(data, values...)
	}

	var result []interface{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := v.(error); ok {
			return nil, false, err
		}

		result = append(result, v)
	}

	return result, len(result) != 0, nil
}

func HandleQueryResult(
	result []interface{},
	marshalFunc func(in interface{}) ([]byte, error),
) string {
	switch len(result) {
	case 0:
		return ""
	case 1:
		switch r := result[0].(type) {
		case string:
			return r
		case []byte:
			return string(r)
		case []interface{}, map[string]interface{}:
			res, _ := marshalFunc(r)
			return string(res)
		case int64:
			return strconv.FormatInt(r, 10)
		case int32:
			return strconv.FormatInt(int64(r), 10)
		case int16:
			return strconv.FormatInt(int64(r), 10)
		case int8:
			return strconv.FormatInt(int64(r), 10)
		case int:
			return strconv.FormatInt(int64(r), 10)
		case uint64:
			return strconv.FormatUint(r, 10)
		case uint32:
			return strconv.FormatUint(uint64(r), 10)
		case uint16:
			return strconv.FormatUint(uint64(r), 10)
		case uint8:
			return strconv.FormatUint(uint64(r), 10)
		case uint:
			return strconv.FormatUint(uint64(r), 10)
		case float64:
			return strconv.FormatFloat(r, 'f', -1, 64)
		case float32:
			return strconv.FormatFloat(float64(r), 'f', -1, 64)
		case bool:
			return strconv.FormatBool(r)
		case nil:
			return "null"
		default:
			return fmt.Sprintf("%v", r)
		}
	default:
		res, _ := marshalFunc(result)
		return string(res)
	}
}
