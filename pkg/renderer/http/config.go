package http

import (
	"net/http"

	"arhat.dev/dukkha/pkg/field"
)

type rendererHTTPConfig struct {
	field.BaseField

	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type inputHTTPConfig struct {
	field.BaseField

	URL string `yaml:"url"`

	Config rendererHTTPConfig `yaml:",inline"`
}

func (c *rendererHTTPConfig) createClient() *http.Client {
	// TODO: create client according to config value
	return http.DefaultClient
}
