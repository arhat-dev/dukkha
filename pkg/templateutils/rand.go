package templateutils

// TODO: use crypto/rand as math/rand
// nolint:unused
type randNS struct{}

// Int(): return an integer in range [0, 100]
// Int(max Number): return an integer in range [0, max]
// Int(min, max Number): return an integer in range [min, max]
// nolint:unused
func (randNS) Int(args ...Number) (_ int64, err error) {
	return
}

// Float(): return a float in range [0,1]
// Float(max Number): return a float in range [0, max]
// Float(min, max Number): return a float in range [min, max]
// nolint:unused
func (randNS) Float(args ...Number) (_ float64, err error) {
	return
}
