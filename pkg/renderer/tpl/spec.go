package tpl

import "arhat.dev/rs"

type configSpec struct {
	// Include templates from local files/dirs
	Include []string `yaml:"include"`

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
