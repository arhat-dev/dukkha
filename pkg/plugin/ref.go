package plugin

import "arhat.dev/rs"

type PluginReference struct {
	rs.BaseField

	Name string `yaml:"name"`
}

type PluginSource struct {
	rs.BaseField

	// Source of the single file plugin go code
	Source *string `yaml:"source"`

	// Module name of plugin, MUST have its dependencies vendored
	Module *string `yaml:"module"`
}
