package matrix

func NewFilter(match map[string][]string) *Filter {
	mv := make(map[string]*Vector, len(match))
	for k, v := range match {
		mv[k] = NewVector(v...)
	}

	return &Filter{
		match:  mv,
		ignore: nil,
	}
}

type Filter struct {
	match  map[string]*Vector
	ignore [][2]string
}

func (f *Filter) Equals(a *Filter) bool {
	if f == nil {
		if a == nil {
			return true
		}

		return false
	}

	if a == nil {
		return false
	}

	if len(f.match) != len(a.match) || len(f.ignore) != len(a.ignore) {
		return false
	}

	for k, v := range f.match {
		va, ok := a.match[k]
		if !ok {
			return false
		}

		if !v.Equals(va) {
			return false
		}
	}

	for i, v := range f.ignore {
		va := a.ignore[i]

		if v[0] != va[0] || v[1] != va[1] {
			return false
		}
	}

	return true
}

func (f *Filter) AddMatch(key, value string) {
	vec, ok := f.match[key]
	if ok {
		vec.Vector = append(vec.Vector, value)
	} else {
		f.match[key] = NewVector(value)
	}
}

func (f *Filter) AddIgnore(key, value string) {
	f.ignore = append(f.ignore, [2]string{key, value})
}

// AsEntry converts f.match to a matrix Entry (used for task matrix)
// should only be used when you are sure the matrix filter is set
// for your task matrix execution
func (f *Filter) AsEntry() Entry {
	if f == nil {
		return nil
	}

	ret := make(map[string]string, len(f.match))
	for k, v := range f.match {
		if len(v.Vector) == 0 {
			ret[k] = ""
		} else {
			ret[k] = v.Vector[0]
		}
	}

	return ret
}

func (f *Filter) Clone() *Filter {
	var (
		matchFilter  = make(map[string]*Vector, len(f.match))
		ignoreFilter = make([][2]string, len(f.ignore))
	)

	for k, v := range f.match {
		matchFilter[k] = NewVector(v.Vector...)
	}

	for i, kv := range f.ignore {
		ignoreFilter[i] = [2]string{kv[0], kv[1]}
	}

	return &Filter{
		match:  matchFilter,
		ignore: ignoreFilter,
	}
}
