package matrix

// A Filter represents a set of match/ignore rules of a matrix
type Filter struct {
	match  map[string]*Vector
	ignore [][2]string
}

// Equals returns true when everything in x is the same with what in f
func (f *Filter) Equals(x *Filter) bool {
	if f == nil {
		return x == nil
	}

	if x == nil {
		return false
	}

	if len(f.match) != len(x.match) || len(f.ignore) != len(x.ignore) {
		return false
	}

	for k, v := range f.match {
		va, ok := x.match[k]
		if !ok {
			return false
		}

		if !v.Equals(va) {
			return false
		}
	}

	for i, v := range f.ignore {
		va := x.ignore[i]

		if v[0] != va[0] || v[1] != va[1] {
			return false
		}
	}

	return true
}

// AddMatch adds a key value match pair
func (f *Filter) AddMatch(key, value string) {
	if f.match == nil {
		f.match = make(map[string]*Vector)
	}

	vec, ok := f.match[key]
	if ok {
		vec.Vec = append(vec.Vec, value)
	} else {
		f.match[key] = NewVector(value)
	}
}

// AddIgnore adds a ignore rule matching key=value
func (f *Filter) AddIgnore(key, value string) {
	f.ignore = append(f.ignore, [2]string{key, value})
}

// AsEntry converts f.match to a matrix [Entry] (used for task matrix)
//
// should only be used when you are sure the matrix filter is set
// for your task matrix execution
func (f *Filter) AsEntry() Entry {
	if f == nil {
		return nil
	}

	ret := make(map[string]string, len(f.match))
	for k, v := range f.match {
		if len(v.Vec) == 0 {
			ret[k] = ""
		} else {
			ret[k] = v.Vec[0]
		}
	}

	return ret
}

// Clone deep copies f to a new filter
func (f *Filter) Clone() (ret Filter) {
	ret.match = make(map[string]*Vector, len(f.match))
	ret.ignore = make([][2]string, len(f.ignore))

	for k, v := range f.match {
		ret.match[k] = NewVector(v.Vec...)
	}

	for i, kv := range f.ignore {
		ret.ignore[i] = [2]string{kv[0], kv[1]}
	}

	return
}

// IsEmpty returns true when there is nothing in f
func (f *Filter) IsEmpty() bool {
	return f == nil || (len(f.match) == 0 && len(f.ignore) == 0)
}
