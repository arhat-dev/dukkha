package templateutils

import (
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
)

type uuidNS struct{}

const (
	zeroUUID = "00000000-0000-0000-0000-000000000000"
)

func uuidv1() (string, error) {
	u, err := uuid.NewUUID()
	return u.String(), err
}

func uuidv4() (string, error) {
	u, err := uuid.NewRandomFromReader(rand.Reader)
	return u.String(), err
}

// Zero returns a uuid with all bits set to zero
func (uuidNS) Zero() string        { return zeroUUID }
func (uuidNS) V1() (string, error) { return uuidv1() }
func (uuidNS) V4() (string, error) { return uuidv4() }

func (uuidNS) IsValid(in String) bool {
	s, err := toString(in)
	if err != nil {
		return false
	}

	_, err = uuid.Parse(s)
	return err == nil
}

// New return a uuid by optional version arg
// when there is no version arg, return v4 uuid
func (uuidNS) New(version ...Number) (string, error) {
	switch n := len(version); n {
	case 0:
		return uuidv4()
	default:
		ver, err := parseInteger[uint64](version[n-1])
		if err != nil {
			return "", err
		}

		switch ver {
		case 1:
			return uuidv1()
		case 4:
			return uuidv4()
		default:
			return "", fmt.Errorf("unsupported uuid version %d", ver)
		}
	}
}
