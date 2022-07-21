package render

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"arhat.dev/rs"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
)

func renderYamlReader(
	rc dukkha.Context,
	src io.Reader,
	destPath *string,
	destPerm fs.FileMode,
	opts *ResolvedOptions,
) error {
	dec := yaml.NewDecoder(src)
	ret, err := parseYaml(dec)

	// always write parsed yaml docs, regardless of errors

	var enc encoder
	if len(ret) != 0 {
		if destPath == nil {
			enc, err = opts.createEncoder(nil)
			if err != nil {
				return err
			}
		} else {
			var _destFile fs.File
			_destFile, err = rc.FS().OpenFile(
				*destPath,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				destPerm,
			)
			if err != nil {
				return fmt.Errorf("open output file %q: %w",
					*destPath, err,
				)
			}
			destFile := _destFile.(*os.File)
			defer func() { _ = destFile.Close() }()

			enc, err = opts.createEncoder(destFile)
			if err != nil {
				return err
			}
		}
	}

	for _, doc := range ret {
		err2 := doc.ResolveFields(rc, -1)
		if err2 != nil {
			return multierr.Append(err, err2)
		}

		err2 = enc.Encode(doc)
		if err2 != nil {
			return multierr.Append(err, err2)
		}
	}

	return err
}

func parseYaml(dec *yaml.Decoder) ([]*rs.AnyObject, error) {
	var ret []*rs.AnyObject
	for {
		obj := &rs.AnyObject{}
		err := dec.Decode(obj)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return ret, nil
			}

			return ret, err
		}

		ret = append(ret, obj)
	}
}
