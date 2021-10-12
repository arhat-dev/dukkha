package textquery

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
)

// Query runs jq query over general text data bytes with custom
// marshaling/unmarshaling func for data serialization/deserialization
func Query(
	query string,
	generatNext func() (interface{}, bool),
	marshalFunc func(in interface{}) ([]byte, error),
) (string, error) {
	q, err := gojq.Parse(query)
	if err != nil {
		return "", fmt.Errorf("failed to parse query: %w", err)
	}

	sb := &strings.Builder{}
	wroteOnce := false
	for {
		data, ok := generatNext()
		if !ok {
			break
		}

		result, hasResult, err2 := RunQuery(q, data, nil)
		if hasResult {
			if wroteOnce {
				sb.WriteByte('\n')
			}

			sb.WriteString(HandleQueryResult(result, marshalFunc))

			if !wroteOnce {
				wroteOnce = true
			}
		}

		if err2 != nil {
			return sb.String(), err2
		}
	}

	return sb.String(), nil
}

// RunQuery runs jq query over arbitrary data with optional
// predefined key value pairs
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

// HandleQueryResult from RunQuery
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
