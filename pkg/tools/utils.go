package tools

import (
	"reflect"
	"strings"
)

func getTagNamesToResolve(typ reflect.Type) []string {
	var ret []string
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if len(f.PkgPath) != 0 {
			continue
		}

		switch f.Name {
		case "BaseField":
			continue
		}

		yTags := strings.Split(f.Tag.Get("yaml"), ",")
		if yTags[0] == "-" {
			// ignored explicitly
			continue
		}

		isInline := false
		for _, tag := range yTags[1:] {
			switch tag {
			case "inline":
				isInline = true
				// inline field can only be struct or map
				if f.Type.Kind() == reflect.Map {
					ret = append(ret, f.Name)
				} else {
					ret = append(ret, getTagNamesToResolve(f.Type)...)
				}
			default:
			}
		}

		if isInline {
			continue
		}

		tagName := yTags[0]
		if len(tagName) == 0 {
			tagName = strings.ToLower(f.Name)
		}

		ret = append(ret, tagName)
	}

	return ret
}
