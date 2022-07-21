package parse

import "reflect"

// TemplateFuncs defines methods required to become a template func provider
type TemplateFuncs interface {
	// Has returns true when there is template func with the same name
	Has(name string) bool

	// GetByName returns the template func with the same name
	GetByName(name string) reflect.Value
}
