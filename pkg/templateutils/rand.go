package templateutils

import (
	"math/rand"
	"time"
)

type randNS struct{}

// Int(): return an integer in range [0, 100]
// Int(max Number): return an integer in range [0, max]
// Int(min, max Number): return an integer in range [min, max]
func (randNS) Int(args ...Number) (_ int64, err error) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return
}

// Float(): return a float in range [0,1]
// Float(max Number): return a float in range [0, max]
// Float(min, max Number): return a float in range [min, max]
func (randNS) Float(args ...Number) (_ float64, err error) {
	return
}
