package render

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"arhat.dev/pkg/textquery"
	"arhat.dev/rs"
	"github.com/itchyny/gojq"
	"gopkg.in/yaml.v3"
)

type (
	encoder interface {
		Encode(v *rs.AnyObject) error
	}
	encoderCreateFunc func(w io.Writer) (encoder, error)
)

var _ encoder = (*encoderWithQueryRequest)(nil)

type encoderWithQueryRequest struct {
	encode func(interface{}) error
	query  *gojq.Query
}

func (e *encoderWithQueryRequest) Encode(v *rs.AnyObject) error {
	if e.query != nil {
		ret, _, err := textquery.RunQuery(e.query, v.NormalizedValue(), nil)
		if err != nil {
			return err
		}

		switch len(ret) {
		case 0:
			return e.encode(nil)
		case 1:
			return e.encode(ret[0])
		default:
			return e.encode(ret)
		}
	} else {
		return e.encode(v)
	}
}

func newEncoder(
	query *gojq.Query,
	w io.Writer,
	outputFormat, indentStr string,
	indentSize int,
) (encoder, error) {

	var encImpl func(interface{}) error
	switch outputFormat {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", strings.Repeat(indentStr, indentSize))
		encImpl = enc.Encode
	case "yaml":
		fallthrough
	case "":
		enc := yaml.NewEncoder(w)
		enc.SetIndent(indentSize)
		encImpl = enc.Encode
	default:
		return nil, fmt.Errorf("unknown output format %q", outputFormat)
	}

	return &encoderWithQueryRequest{
		encode: encImpl,
		query:  query,
	}, nil
}
