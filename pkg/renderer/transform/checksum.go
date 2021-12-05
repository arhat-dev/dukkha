package transform

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/rs"
)

type checksumKind string

// nolint:revive
const (
	checksum_MD5    checksumKind = "md5"
	checksum_SHA1   checksumKind = "sha1"
	checksum_SHA224 checksumKind = "sha224"
	checksum_SHA256 checksumKind = "sha256"
	checksum_SHA512 checksumKind = "sha512"
)

type Checksum struct {
	rs.BaseField `yaml:"-" json:"-"`

	// File is the local file path to the target data file
	//
	// File and Data are mutually exclusive
	File *string `yaml:"file"`

	// Data is the raw string to be verified
	//
	// Path and Data are mutually exclusive
	Data *string `yaml:"data"`

	// Kind is the name of hash algo
	Kind checksumKind `yaml:"kind"`

	// Sum is the expected hex encoded string of checksum value
	Sum string `yaml:"sum"`

	// Key is the optional hmac key, if set, will do hmac
	Key *string `yaml:"key"`
}

// VerifyFile check local file if matching the checksum
func (cs Checksum) Verify(ofs *fshelper.OSFS) error {
	var newHash func() hash.Hash
	switch cs.Kind {
	case checksum_MD5:
		newHash = md5.New
	case checksum_SHA1:
		newHash = sha1.New
	case checksum_SHA224:
		newHash = sha256.New224
	case checksum_SHA256:
		newHash = sha256.New
	case checksum_SHA512:
		newHash = sha512.New
	default:
		return fmt.Errorf("unknown checksum type %q", cs.Kind)
	}

	var h hash.Hash
	if cs.Key != nil {
		h = hmac.New(newHash, bytes.TrimSpace([]byte(*cs.Key)))
	} else {
		h = newHash()
	}

	var src io.Reader
	switch {
	case cs.File != nil:
		f, err := ofs.Open(*cs.File)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		src = f
	case cs.Data != nil:
		src = bytes.NewReader([]byte(*cs.Data))
	default:
		return fmt.Errorf("nothing to check")
	}

	_, err := io.Copy(h, src)
	if err != nil {
		return err
	}

	sum := hex.EncodeToString(h.Sum(nil))
	if sum == cs.Sum {
		return nil
	}

	return fmt.Errorf(
		"data corrupted: checksum (%s:%s) not match",
		cs.Kind, sum,
	)
}
