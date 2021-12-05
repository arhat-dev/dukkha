package conf

import (
	"arhat.dev/rs"
)

type IncludeEntry struct {
	rs.BaseField `yaml:",inline"`

	// Path is the local path to include, can be either directory or file
	//
	// Path and Text are mutually exclusive
	Path string `yaml:"path"`

	// Text is the config text to include, usually used with rendering suffix
	// to include remote config
	//
	// Path and Text are mutually exclusive
	Text string `yaml:"text"`
}
