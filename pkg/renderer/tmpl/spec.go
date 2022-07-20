package tmpl

import "arhat.dev/rs"

type includeSpec struct {
	rs.BaseField

	// Path to local files/dirs
	Path string `yaml:"path"`

	// Text is the plain text template to be included
	Text string `yaml:"text"`
}

type configSpec struct {
	rs.BaseField

	// Include templates
	Include []*includeSpec `yaml:"include"`

	// Variables are a map of any data
	//
	// available as `var.some_value`
	Variables rs.AnyObjectMap `yaml:"variables"`
}

type inputSpec struct {
	rs.BaseField

	// Template text
	Template string `yaml:"template"`

	Config configSpec `yaml:",inline"`
}
