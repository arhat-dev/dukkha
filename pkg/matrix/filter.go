package matrix

func NewFilter(match map[string][]string) *Filter {
	if match == nil {
		match = make(map[string][]string)
	}

	return &Filter{
		match:  match,
		ignore: nil,
	}
}

type Filter struct {
	match  map[string][]string
	ignore [][2]string
}

func (f *Filter) AddMatch(key, value string) {
	f.match[key] = append(f.match[key], value)
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
		if len(v) == 0 {
			ret[k] = ""
		} else {
			ret[k] = v[0]
		}
	}

	return ret
}

func (f *Filter) Clone() *Filter {
	var (
		matchFilter  = make(map[string][]string, len(f.match))
		ignoreFilter = make([][2]string, len(f.ignore))
	)

	for k, v := range f.match {
		matchFilter[k] = append(matchFilter[k], v...)
	}

	for i, kv := range f.ignore {
		ignoreFilter[i] = [2]string{kv[0], kv[1]}
	}

	return &Filter{
		match:  matchFilter,
		ignore: ignoreFilter,
	}
}
