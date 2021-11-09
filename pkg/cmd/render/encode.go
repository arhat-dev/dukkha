package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"arhat.dev/pkg/textquery"
	"arhat.dev/rs"
	"github.com/itchyny/gojq"
	"gopkg.in/yaml.v3"
)

type encoder interface {
	Encode(v *rs.AnyObject) error
}

var _ encoder = (*anyObjectEncoder)(nil)

type anyObjectEncoder struct {
	writeDocSep bool

	write  func([]byte) (int, error)
	encode func(interface{}) error
	query  *gojq.Query
}

func (e *anyObjectEncoder) Encode(v *rs.AnyObject) error {
	var toEncode interface{}
	if e.query != nil {
		ret, _, err := textquery.RunQuery(e.query, v.NormalizedValue(), nil)
		if err != nil {
			return err
		}

		if len(ret) == 0 {
			return e.encode(nil)
		}

		if len(ret) != 1 {
			return e.encode(ret)
		}

		toEncode = ret[0]
	} else {
		toEncode = v.NormalizedValue()
	}

	if e.writeDocSep {
		e.writeDocSep = false
		_, err := e.write([]byte("---\n"))
		if err != nil {
			return err
		}
	}

	// special case []byte and string, which can be produced by virtual key rederer
	switch rt := toEncode.(type) {
	case []byte:
		_, err := e.write(rt)
		if err != nil {
			return err
		}

		if !bytes.HasSuffix(rt, []byte{'\n'}) {
			_, err = e.write([]byte{'\n'})
			if err != nil {
				return err
			}
		}

		e.writeDocSep = true
		return nil
	case string:
		_, err := e.write([]byte(rt))
		if err != nil {
			return err
		}

		if !strings.HasSuffix(rt, "\n") {
			_, err = e.write([]byte{'\n'})
			if err != nil {
				return err
			}
		}

		e.writeDocSep = true
		return nil
	default:
		return e.encode(rt)
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

	return &anyObjectEncoder{
		write:  w.Write,
		encode: encImpl,
		query:  query,
	}, nil
}
