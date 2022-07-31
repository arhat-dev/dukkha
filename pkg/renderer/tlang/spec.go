package tlang

import "arhat.dev/rs"

type IncludeSpec struct {
	rs.BaseField

	// Path to local files/dirs
	Path string `yaml:"path"`

	// Text is the plain text template to be included
	Text string `yaml:"text"`
}

type ConfigSpec struct {
	rs.BaseField

	// Include templates
	Include []*IncludeSpec `yaml:"include"`

	// Variables are a map of any data
	//
	// available as `var.some_value`
	Variables rs.AnyObjectMap `yaml:"variables"`
}

type InputSpec struct {
	rs.BaseField

	// Script text of tlang
	Script string `yaml:"script"`

	Config ConfigSpec `yaml:",inline"`
}
