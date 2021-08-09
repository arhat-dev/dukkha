package http

import (
	"net/http"

	"arhat.dev/rs"
)

type rendererHTTPConfig struct {
	rs.BaseField

	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type inputHTTPConfig struct {
	rs.BaseField

	URL string `yaml:"url"`

	Config rendererHTTPConfig `yaml:",inline"`
}

func (c *rendererHTTPConfig) createClient() *http.Client {
	// TODO: create client according to config value
	return http.DefaultClient
}
