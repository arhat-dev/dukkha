package templateutils

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode"

	"arhat.dev/pkg/clihelper"
	"arhat.dev/pkg/stringhelper"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// encNS for encoding
type encNS struct{}

// YAML encodes object into yaml
//
// NOTE: there's no yaml decoding support in this namespace, use dukkha.FromYaml instead
func (encNS) YAML(args ...any) (ret string, err error) {
	obj, outWriter, args := parseArgs_MaybeOUTPUT_OBJ(args)
	flags, err := toStrings(args)
	if err != nil {
		return
	}

	doEncode := func(indentN int) {
		var (
			tmpWriter bytes.Buffer
			enc       *yaml.Encoder
		)

		if outWriter != nil {
			enc = yaml.NewEncoder(outWriter)
		} else {
			enc = yaml.NewEncoder(&tmpWriter)
		}

		enc.SetIndent(indentN)

		err = enc.Encode(obj)
		ret = stringhelper.Convert[string, byte](tmpWriter.Next(tmpWriter.Len()))
	}

	if len(flags) == 0 {
		doEncode(2)
		return
	}

	var (
		fs pflag.FlagSet

		indentN int
	)

	clihelper.InitFlagSet(&fs, "yaml")

	fs.IntVarP(&indentN, "indent-count", "c", 2, "") // keep it compatible with JSON

	err = fs.Parse(flags)
	if err != nil {
		return
	}

	doEncode(indentN)
	return
}

// Json to encode object to json format
// Usage: Json(...<options>, <optional writer>, Object)
//
// NOTE: there's no json decoding support in this namespace, use dukkha.FromJson instead
func (encNS) JSON(args ...any) (ret string, err error) {
	obj, outWriter, args := parseArgs_MaybeOUTPUT_OBJ(args)
	flags, err := toStrings(args)
	if err != nil {
		return
	}

	doEncode := func(indentN int, indentStr string, escapeHTML bool) {
		var (
			tmpWriter bytes.Buffer
			enc       *json.Encoder
		)

		if outWriter != nil {
			enc = json.NewEncoder(outWriter)
		} else {
			enc = json.NewEncoder(&tmpWriter)
		}

		enc.SetIndent("", strings.Repeat(indentStr, indentN))
		enc.SetEscapeHTML(escapeHTML)

		err = enc.Encode(obj)
		ret = stringhelper.Convert[string, byte](tmpWriter.Next(tmpWriter.Len()))
	}

	if len(flags) == 0 {
		doEncode(0, "", true)
		return
	}

	var (
		fs pflag.FlagSet

		indentN    int
		indentStr  string
		escapeHTML bool
		pretty     bool
	)

	clihelper.InitFlagSet(&fs, "json")

	fs.IntVarP(&indentN, "indent-count", "c", 2, "") // keep it compatible with YAML
	fs.StringVarP(&indentStr, "indent", "i", "", "")
	fs.BoolVarP(&escapeHTML, "escape-html", "e", true, "")
	fs.BoolVarP(&pretty, "pretty", "P", false, "")

	err = fs.Parse(flags)
	if err != nil {
		return
	}

	if pretty {
		if len(indentStr) == 0 {
			indentN, indentStr = 2, " "
		}
	}

	doEncode(indentN, indentStr, escapeHTML)
	return
}

// Hex to encode/decode hex data (default works in encoding mode)
//
// Hex(...<options>, <optional writer>, inputDataOrReader)
//
// where options are:
// - `--decode` or `-d`: decode input as hex
func (encNS) Hex(args ...any) (ret string, err error) {
	inData, inReader, outWriter, args, err := parseArgs_MaybeOUTPUT_DATA(args)
	if err != nil {
		return
	}

	flags, err := toStrings(args)
	if err != nil {
		return
	}

	// called at most once
	doEncode := func() {
		var buf []byte
		buf, err = writeBytes(inData, inReader, outWriter,
			/* pre-check */ func(tmpReader *bytes.Reader, tmpWriter *bytes.Buffer) io.Writer {
				return nil
			},
			/* before writing */ hex.NewEncoder,
			/* after wrote */ func(wrappedDst io.Writer) {},
			/* fallback */ func(b []byte) ([]byte, error) {
				return stringhelper.ToBytes[byte, byte](hex.EncodeToString(b)), nil
			},
		)

		ret = stringhelper.Convert[string, byte](buf)
	}

	if len(flags) == 0 {
		doEncode()
		return
	}

	var (
		fs pflag.FlagSet

		decode bool
	)

	clihelper.InitFlagSet(&fs, "hex")

	fs.BoolVarP(&decode, "decode", "d", false, "")

	err = fs.Parse(flags)
	if err != nil {
		return
	}

	if decode {
		var buf []byte

		buf, err = readBytes(inData, inReader, outWriter,
			/* pre-check */ func(tmpReader *bytes.Reader, tmpWriter *bytes.Buffer) io.Reader {
				return nil
			},
			/* before reading */ hex.NewDecoder,
			/* fallback */ func(b []byte) ([]byte, error) {
				return hex.DecodeString(stringhelper.Convert[string, byte](b))
			},
		)
		if err != nil {
			return
		}

		return stringhelper.Convert[string, byte](buf), nil
	}

	doEncode()
	return
}

type iEncoding interface {
	EncodeToString(src []byte) string
	DecodeString(s string) ([]byte, error)
}

type baseXFactory struct {
	Std      iEncoding // required
	URL, Hex iEncoding // optional

	ToRawEncoding    func(iEncoding) iEncoding                 // required
	ToStrictEncoding func(iEncoding) iEncoding                 // required
	NewEncoding      func(string) iEncoding                    // required
	CreateEncoder    func(iEncoding, io.Writer) io.WriteCloser // required
	CreateDecoder    func(iEncoding, io.Reader) io.Reader      // required
}

// Base64 to encode/decode data in base64 format
// Usage: Base64(...<options>, <optional writer>, inputDataOrReader)
// where options are:
// - `--decode` or `-d`: decode input as base64 encoded data
// - `--wrap` or `-w` N: (only for encoding) wrap each line into N bytes (except last line)
// - `--url` or `-u`: use URL encoding (for base64)
// - `--hex` or `-h`: use Hex encoding (for base32)
// - `--raw` or `-r`: encode/decode without padding
// - `--strict` or `-s`: encode/decode in strict mode
// - `--table` or `-t` str: encode/decode with custom encoding
// 							table length MUST match (e.g. 32 for base32, 64 for base64)
func (encNS) Base64(args ...any) (string, error) {
	const (
		encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
		encodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	)

	return handleBaseX(args, "base64", baseXFactory{
		Std: base64.NewEncoding(encodeStd),
		URL: base64.NewEncoding(encodeURL),

		ToRawEncoding: func(e iEncoding) iEncoding {
			return e.(*base64.Encoding).WithPadding(base64.NoPadding)
		},
		ToStrictEncoding: func(in iEncoding) iEncoding {
			return in.(*base64.Encoding).Strict()
		},
		NewEncoding: func(s string) iEncoding {
			return base64.NewEncoding(s)
		},
		CreateEncoder: func(e iEncoding, w io.Writer) io.WriteCloser {
			return base64.NewEncoder(e.(*base64.Encoding), w)
		},
		CreateDecoder: func(e iEncoding, r io.Reader) io.Reader {
			return base64.NewDecoder(e.(*base64.Encoding), r)
		},
	})
}

// Base32 see comments for Base64 (replace 64 with 32)
func (encNS) Base32(args ...any) (string, error) {
	const (
		encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
		encodeHex = "0123456789ABCDEFGHIJKLMNOPQRSTUV"
	)

	return handleBaseX(args, "base32", baseXFactory{
		Std: base32.NewEncoding(encodeStd),
		Hex: base32.NewEncoding(encodeHex),

		ToStrictEncoding: func(ie iEncoding) iEncoding {
			return nil
		},
		NewEncoding: func(s string) iEncoding {
			return base32.NewEncoding(s)
		},
		CreateEncoder: func(e iEncoding, w io.Writer) io.WriteCloser {
			return base32.NewEncoder(e.(*base32.Encoding), w)
		},
		CreateDecoder: func(e iEncoding, r io.Reader) io.Reader {
			return base32.NewDecoder(e.(*base32.Encoding), r)
		},
		ToRawEncoding: func(e iEncoding) iEncoding {
			return e.(*base32.Encoding).WithPadding(base32.NoPadding)
		},
	})
}

func handleBaseX(
	args []any,
	name string,
	factory baseXFactory,
) (ret string, err error) {
	inData, inReader, outWriter, args, err := parseArgs_MaybeOUTPUT_DATA(args)
	if err != nil {
		return
	}

	flags, err := toStrings(args)
	if err != nil {
		return
	}

	// called at most once
	doEncode := func(enc iEncoding, lineWrap int) {
		var (
			buf []byte
			cw  ChunkedWriter
		)
		buf, err = writeBytes(inData, inReader, outWriter,
			/* pre-check */ func(tmpReader *bytes.Reader, tmpWriter *bytes.Buffer) io.Writer {
				if lineWrap > 0 {
					underlay := outWriter
					if underlay == nil {
						underlay = tmpWriter
					}

					cw = NewChunkedWriter(
						lineWrap,
						underlay,
						func() error { return nil },
						func() error {
							_, err = underlay.Write([]byte("\n"))
							return err
						},
					)

					return &cw
				}

				return nil
			},
			/* before writing */ func(w io.Writer) io.Writer {
				return factory.CreateEncoder(enc, w)
			},
			/* after wrote */ func(wrappedDst io.Writer) {
				_ = wrappedDst.(io.Closer).Close()
			},
			/* fallback */ func(b []byte) ([]byte, error) {
				str := enc.EncodeToString(b)
				return []byte(str), nil
			},
		)

		if len(buf) != 0 {
			ret = stringhelper.Convert[string, byte](buf)
		}
	}

	if len(flags) == 0 {
		doEncode(factory.Std, 0)
		return
	}

	// keep FlagSet on stack
	var (
		fs pflag.FlagSet

		strict bool
		decode bool

		url bool
		raw bool
		hex bool

		table string
		wrap  int
	)

	clihelper.InitFlagSet(&fs, name)

	// decode/encode
	fs.BoolVarP(&decode, "decode", "d", false, "")

	// control encoding
	fs.BoolVarP(&strict, "strict", "s", false, "") // for base64
	fs.BoolVarP(&url, "url", "u", false, "")       // for base64
	fs.BoolVarP(&hex, "hex", "h", false, "")       // for base32
	fs.BoolVarP(&raw, "raw", "r", false, "")       // for base32, base64
	fs.StringVarP(&table, "table", "t", "", "")    // for base32, base64

	// control line wrapping width (default gnu/linux defaults to 76, we don't do that)
	// only applicapable to encode mode
	fs.IntVarP(&wrap, "wrap", "w", 0, "")

	err = fs.Parse(flags)
	if err != nil {
		return
	}

	var targetEncoding iEncoding

	switch {
	case len(table) != 0:
		targetEncoding = factory.NewEncoding(table)
	case !url && !hex:
		targetEncoding = factory.Std
	case url && !hex:
		targetEncoding = factory.URL
	case hex && !url:
		targetEncoding = factory.Hex
	}

	if targetEncoding == nil {
		err = fmt.Errorf("invalid encoding options: url = %v, raw = %v, hex = %v", url, raw, hex)
		return
	}

	if raw {
		targetEncoding = factory.ToRawEncoding(targetEncoding)
		if targetEncoding == nil {
			err = fmt.Errorf("raw encoding not supported")
			return
		}
	}

	if strict {
		targetEncoding = factory.ToStrictEncoding(targetEncoding)
		if targetEncoding == nil {
			err = fmt.Errorf("strict mode not supported")
			return
		}
	}

	if decode {
		var buf []byte
		buf, err = readBytes(inData, inReader, outWriter,
			/* pre-check */ func(tmpReader *bytes.Reader, tmpWriter *bytes.Buffer) io.Reader {
				if inReader == nil {
					tmpReader.Reset(inData)
					inReader = tmpReader
				}

				return NewFilterReader(inReader, func(p []byte) int {
					return RemoveMatchedRunesInPlace(p, unicode.IsSpace)
				})
			},
			/* before read */ func(r io.Reader) io.Reader {
				return factory.CreateDecoder(targetEncoding, r)
			},
			/* fallback */ func(b []byte) ([]byte, error) {
				return targetEncoding.DecodeString(stringhelper.Convert[string, byte](b))
			},
		)

		if len(buf) != 0 {
			ret = stringhelper.Convert[string, byte](buf)
		}

		return
	}

	doEncode(targetEncoding, wrap)
	return
}

func readBytes(
	inData []byte,
	inReader io.Reader,
	outWriter io.Writer,

	// called before doing anything
	preCheck func(tmpReader *bytes.Reader, tmpWriter *bytes.Buffer) io.Reader,

	// called before writing to dst, returned writer will replace dst (MUST NOT be nil)
	beforeRead func(src io.Reader) io.Reader,

	// called when there is no reader/writer to handle inData
	handleFallback func([]byte) ([]byte, error),
) (ret []byte, err error) {
	var (
		tmpReader bytes.Reader
		tmpWriter bytes.Buffer
	)

	x := preCheck(&tmpReader, &tmpWriter)
	if x != nil {
		inReader = x
	}

	// nolint:gocritic
	if outWriter != nil {
		if inReader == nil {
			tmpReader.Reset(inData)
			inReader = &tmpReader
		}

		_, err = io.Copy(outWriter, beforeRead(inReader))
		ret = tmpWriter.Next(tmpWriter.Len())
		return
	} else if inReader != nil {
		if outWriter == nil {
			outWriter = &tmpWriter
		}

		_, err = io.Copy(outWriter, beforeRead(inReader))
		ret = tmpWriter.Next(tmpWriter.Len())
		return
	} else {
		return handleFallback(inData)
	}
}

func writeBytes(
	inData []byte,
	inReader io.Reader,
	outWriter io.Writer,

	// called before doing anything
	preCheck func(tmpReader *bytes.Reader, tmpWriter *bytes.Buffer) io.Writer,

	// called before writing to dst, returned writer will replace dst (MUST NOT be nil)
	beforeWriting func(dst io.Writer) io.Writer,

	// called after finished copying data to writer returned by beforeWrite()
	afterWrote func(wrappedDst io.Writer),

	// called when there is no reader/writer to handle inData
	handleFallback func([]byte) ([]byte, error),
) (ret []byte, err error) {
	var (
		tmpReader bytes.Reader
		tmpWriter bytes.Buffer
	)

	x := preCheck(&tmpReader, &tmpWriter)
	if x != nil {
		outWriter = x
	}

	// nolint:gocritic
	if outWriter != nil {
		if inReader == nil {
			tmpReader.Reset(inData)
			inReader = &tmpReader
		}

		outWriter = beforeWriting(outWriter)
		_, err = io.Copy(outWriter, inReader)
		afterWrote(outWriter)

		ret = tmpWriter.Next(tmpWriter.Len())
		return
	} else if inReader != nil {
		if outWriter == nil {
			outWriter = &tmpWriter
		}

		outWriter = beforeWriting(outWriter)
		_, err = io.Copy(outWriter, inReader)
		afterWrote(outWriter)

		ret = tmpWriter.Next(tmpWriter.Len())
		return
	} else {
		return handleFallback(inData)
	}
}
