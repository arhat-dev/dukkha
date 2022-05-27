package templateutils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"io"
	"strings"

	"arhat.dev/pkg/clihelper"
	"arhat.dev/pkg/stringhelper"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
)

// hashNS for hashing and hmac
// all functions (except Bcrypt) in this namespace shares the same arguments support
// XXX(...<options>, data Bytes)
// where options are:
// - `--hmac` or `-k` key: specify hmac key to generate hash mac
// - `--hex`: encode hashed data in hex format (default)
// - `--raw`: do not encode hashed data
// - `--base64`: encode hashed data in base64 format (std encoding)
// - `--base32`: encode hashed data in base32 format (std encoding)
//
// NOTE: to encode hashed data with more advanced options, use pipeline to pass returned hashing result in raw mode to encNS functions
// e.g. {{- hash.CRC32 "-raw" .SomeData | enc.Base64 "--url" "--raw" -}}
type hashNS struct{}

// ADLER32 big-endian
func (hashNS) ADLER32(args ...any) (string, error) {
	var buf [4]byte
	return handleHashTemplateFunc_DATA(args, "adler32",
		func() hash.Hash { return adler32.New() },
		func(data []byte) []byte {
			binary.BigEndian.PutUint32(buf[:], adler32.Checksum(data))
			return buf[:]
		},
	)
}

// CRC32 (Castagnoli) big-endian
func (hashNS) CRC32(args ...any) (string, error) {
	table := crc32.MakeTable(crc32.Castagnoli)
	var buf [4]byte

	return handleHashTemplateFunc_DATA(args, "crc32",
		func() hash.Hash { return crc32.New(table) },
		func(data []byte) []byte {
			binary.BigEndian.PutUint32(buf[:], crc32.Checksum(data, table))
			return buf[:]
		},
	)
}

// CRC64 (ECMA) big-endian
func (hashNS) CRC64(args ...any) (string, error) {
	table := crc64.MakeTable(crc64.ECMA)
	var buf [8]byte

	return handleHashTemplateFunc_DATA(args, "crc64",
		func() hash.Hash { return crc64.New(table) },
		func(data []byte) []byte {
			binary.BigEndian.PutUint64(buf[:], crc64.Checksum(data, table))
			return buf[:]
		},
	)
}

func (hashNS) MD4(args ...any) (_ string, err error) {
	return handleHashTemplateFunc_DATA(args, "md4", md4.New, func(data []byte) []byte {
		h := md4.New()
		_, err = h.Write(data)
		return h.Sum(nil)
	})
}

func (hashNS) MD5(args ...any) (string, error) {
	var buf [md5.Size]byte
	return handleHashTemplateFunc_DATA(args, "md5", md5.New, func(data []byte) []byte {
		buf = md5.Sum(data)
		return buf[:]
	})
}

func (hashNS) RIPEMD160(args ...any) (_ string, err error) {
	return handleHashTemplateFunc_DATA(args, "ripemd160", ripemd160.New, func(data []byte) []byte {
		h := ripemd160.New()
		_, err = h.Write(data)
		return h.Sum(nil)
	})
}

// Bcrypt hashing
//
// Bcrypt(data Bytes): with default cost
//
// Bcrypt(cost Integer, ...<hashing options>, data Bytes): with custom cost
//
// NOTE: `--hmac/-k key` are not supported as hashing options
func (hashNS) Bcrypt(args ...any) (ret string, err error) {
	var buf []byte
	switch n := len(args); n {
	case 0:
		err = errAtLeastOneArgGotZero
		return
	case 1:
		buf, err = toBytes(args[0])
		if err != nil {
			return
		}

		buf, err = bcrypt.GenerateFromPassword(buf, bcrypt.DefaultCost)
		if err != nil {
			return
		}

		ret = hex.EncodeToString(buf)
		return
	default:
		cost := toIntegerOrPanic[int](args[0])
		if cost == 0 {
			cost = bcrypt.DefaultCost
		}

		if cost < bcrypt.MinCost {
			cost = bcrypt.MinCost
		} else if cost > bcrypt.MaxCost {
			cost = bcrypt.MaxCost
		}

		buf, err = toBytes(args[n-1])
		if err != nil {
			return
		}

		buf, err = bcrypt.GenerateFromPassword(buf, cost)
		if err != nil {
			return
		}

		if n == 2 {
			ret = hex.EncodeToString(buf)
			return
		}

		var (
			opts  hashingOptions
			flags []string
		)
		flags, err = toStrings(args[1 : n-1])
		if err != nil {
			return
		}

		opts, err = parseHashingOptions(flags, "bcrypt", true)
		ret = opts.encodeFn(buf)
		return
	}
}

func (hashNS) SHA1(args ...any) (string, error) {
	var buf [sha1.Size]byte
	return handleHashTemplateFunc_DATA(args, "sha1", sha1.New, func(data []byte) []byte {
		buf = sha1.Sum(data)
		return buf[:]
	})
}

func (hashNS) SHA224(args ...any) (string, error) {
	var buf [sha256.Size224]byte
	return handleHashTemplateFunc_DATA(args, "sha224", sha256.New224, func(data []byte) []byte {
		buf = sha256.Sum224(data)
		return buf[:]
	})
}

func (hashNS) SHA256(args ...any) (string, error) {
	var buf [sha256.Size]byte
	return handleHashTemplateFunc_DATA(args, "sha256", sha256.New, func(data []byte) []byte {
		buf = sha256.Sum256(data)
		return buf[:]
	})
}

func (hashNS) SHA384(args ...any) (string, error) {
	var buf [sha512.Size384]byte
	return handleHashTemplateFunc_DATA(args, "sha384", sha512.New384, func(data []byte) []byte {
		buf = sha512.Sum384(data)
		return buf[:]
	})
}

func (hashNS) SHA512(args ...any) (string, error) {
	var buf [sha512.Size]byte
	return handleHashTemplateFunc_DATA(args, "sha512", sha512.New, func(data []byte) []byte {
		buf = sha512.Sum512(data)
		return buf[:]
	})
}

func (hashNS) SHA512_224(args ...any) (string, error) {
	var buf [sha512.Size224]byte
	return handleHashTemplateFunc_DATA(args, "sha512-224", sha512.New512_224, func(data []byte) []byte {
		buf = sha512.Sum512_224(data)
		return buf[:]
	})
}

func (hashNS) SHA512_256(args ...any) (string, error) {
	var buf [sha512.Size256]byte
	return handleHashTemplateFunc_DATA(args, "sha512-256", sha512.New512_256, func(data []byte) []byte {
		buf = sha512.Sum512_256(data)
		return buf[:]
	})
}

type hashingOptions struct {
	encHex    bool
	encBase64 bool
	encBase32 bool
	encRaw    bool

	hmacKey *string
}

func (opts hashingOptions) encodeFn(b []byte) string {
	switch {
	case opts.encBase64:
		return base64.StdEncoding.EncodeToString(b)
	case opts.encBase32:
		return base32.StdEncoding.EncodeToString(b)
	case opts.encRaw:
		return stringhelper.Convert[string, byte](b)
	case opts.encHex:
		return hex.EncodeToString(b)
	default:
		panic("invalid no encoding method set")
	}
}

func parseHashingOptions(flags []string, name string, defaultRaw bool) (ret hashingOptions, err error) {
	var (
		fs pflag.FlagSet

		hmacKey string
	)

	clihelper.InitFlagSet(&fs, name)

	fs.StringVarP(&hmacKey, "hmac", "k", "", "") // hmac key

	fs.BoolVarP(&ret.encHex, "hex", "h", !defaultRaw, "")
	fs.BoolVar(&ret.encBase64, "base64", false, "")
	fs.BoolVar(&ret.encBase32, "base32", false, "")
	fs.BoolVar(&ret.encRaw, "raw", defaultRaw, "")

	for _, f := range flags {
		switch {
		case strings.HasPrefix(f, "--hmac"), strings.HasPrefix(f, "-k"):
			ret.hmacKey = &hmacKey
		}
	}

	err = fs.Parse(flags)
	return
}

func handleHashTemplateFunc_DATA(args []any, name string, newHash func() hash.Hash, sum func(data []byte) []byte) (ret string, err error) {
	n := len(args)
	if n == 0 {
		return "", fmt.Errorf("at least 1 args expected: got 0")
	}

	doHash := func(hmacKey *string, encode func([]byte) string) {
		var (
			h        hash.Hash
			inData   []byte
			inReader io.Reader
			isReader bool
		)
		if hmacKey != nil {
			h = hmac.New(newHash, []byte(*hmacKey))
		}

		inData, inReader, isReader, err = toBytesOrReader(args[n-1])
		if err != nil {
			return
		}

		if isReader {
			if h == nil {
				h = newHash()
			}

			_, err = io.Copy(h, inReader)
			if err != nil {
				return
			}

			ret = encode(h.Sum(nil))
		} else if h != nil {
			_, err = h.Write(inData)
			if err != nil {
				return
			}

			ret = encode(h.Sum(nil))
		} else {
			ret = encode(sum(inData))
		}
	}

	flags, err := toStrings(args[:n-1])
	if err != nil {
		return
	}

	if len(flags) == 0 {
		doHash(nil, hex.EncodeToString)
		return
	}

	opts, err := parseHashingOptions(flags, name, false)
	if err != nil {
		return
	}

	doHash(opts.hmacKey, opts.encodeFn)
	return
}
