package http

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"arhat.dev/pkg/tlshelper"
	"arhat.dev/rs"
	"golang.org/x/net/http/httpproxy"
)

type nameValuePair struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type headers []nameValuePair

type httpProxyConfig struct {
	rs.BaseField `yaml:"-"`

	Enabled bool   `yaml:"enabled"`
	HTTP    string `yaml:"http"`
	HTTPS   string `yaml:"https"`
	NoProxy string `yaml:"no_proxy"`
	CGI     bool   `yaml:"cgi"`
}

type rendererHTTPConfig struct {
	rs.BaseField `yaml:"-"`

	User     string `yaml:"user"`
	Password string `yaml:"password"`

	Headers headers `yaml:"headers"`

	Method string           `yaml:"method"`
	Proxy  *httpProxyConfig `yaml:"proxy"`

	TLS tlshelper.TLSConfig `yaml:"tls"`

	Body *string `yaml:"body"`
}

// inputHTTPSpec for renderer value
type inputHTTPSpec struct {
	rs.BaseField `yaml:"-"`

	URL string `yaml:"url"`

	// Config of the renderer input
	//
	// TODO: we cannot use `yaml:",inline"` due to https://github.com/go-yaml/yaml/issues/362
	// 		 if use a custom type with UnmarshalYAML implmeneted as inline field
	// 		 and the parent field has UnmarshalYAML implementation as well, go-yaml cannot
	// 		 marshal it correctly, which is the case with rs.BaseField
	Config rendererHTTPConfig `yaml:"config"`
}

func (c rendererHTTPConfig) createClient() (*http.Client, error) {
	tlsConfig, err := c.TLS.GetTLSConfig(false)
	if err != nil {
		return nil, err
	}

	var proxy func(req *http.Request) (*url.URL, error)
	if p := c.Proxy; p != nil {
		if p.Enabled {
			cfg := httpproxy.Config{
				HTTPProxy:  p.HTTP,
				HTTPSProxy: p.HTTPS,
				NoProxy:    p.NoProxy,
				CGI:        p.CGI,
			}

			pf := cfg.ProxyFunc()

			proxy = func(req *http.Request) (*url.URL, error) {
				return pf(req.URL)
			}
		}
	} else {
		proxy = http.ProxyFromEnvironment
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
			DialContext: (&net.Dialer{
				// TODO: allow setting timers?
				Timeout:       30 * time.Second,
				KeepAlive:     30 * time.Second,
				FallbackDelay: 300 * time.Millisecond,
			}).DialContext,
			ForceAttemptHTTP2:     tlsConfig != nil,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,

			DialTLSContext:         nil,
			DisableKeepAlives:      false,
			DisableCompression:     false,
			MaxConnsPerHost:        0,
			ResponseHeaderTimeout:  0,
			TLSNextProto:           nil,
			ProxyConnectHeader:     nil,
			MaxResponseHeaderBytes: 0,
			WriteBufferSize:        0,
			ReadBufferSize:         0,
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	return client, nil
}
