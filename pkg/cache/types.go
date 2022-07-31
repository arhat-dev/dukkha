package cache

import "unicode/utf8"

type IdentifiableObject interface {
	// ScopeUniqueID returns a collision free id of the object in owner's scope
	//
	// for example:
	// 		`file` renderer can use file path as scope unique id since local file paths are unique
	//		`http` renderer can use url and info from request/response header to form a scope unique id
	ScopeUniqueID() string

	// Ext to provide extension name of the file, when present
	// the return value should prefixed with dot (.)
	//
	// this is required as some cli app (e.g. golangci-lint) expects a certain
	// file extension, and won't work with files with wrong/no extension name
	Ext() string
}

// IdentifiableString makes plain old strings as IdentifiableObject
type IdentifiableString string

func (s IdentifiableString) ScopeUniqueID() string {
	return string(s)
}

// Ext returns the file extension of s
func (s IdentifiableString) Ext() (ret string) {
	end := len(s)
	for {
		r, sz := utf8.DecodeLastRuneInString(string(s[:end]))
		if r == utf8.RuneError {
			return
		}
		switch r {
		case utf8.RuneError, '/', '\\':
			return
		case '.':
			return string(s[end-sz:])
		}

		end -= sz
	}
}
