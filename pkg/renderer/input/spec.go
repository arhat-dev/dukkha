package input

import "arhat.dev/rs"

type configSpec struct {
	rs.BaseField

	// HideInput do not echo input
	//
	// Defaults to `false`
	HideInput *bool `yaml:"hide_input"`

	// Prompt for user input
	//
	// Defaults to "" (empty)
	Prompt string `yaml:"prompt"`
}

type inputSpec struct {
	rs.BaseField

	Config configSpec `yaml:",inline"`
}
