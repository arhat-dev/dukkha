package kubectl

import "arhat.dev/rs"

type kubeContext struct {
	rs.BaseField `yaml:"-"`

	// KubeconfigFile is the content of your kubeconfig
	Kubeconfig string `yaml:"kubeconfig"`

	// Name of the context (optional if you set required information here instead of using kubeconfig)
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`

	User    string `yaml:"user"`
	Group   string `yaml:"group"`
	Cluster string `yaml:"cluster"`

	APIServer string `yaml:"api_server"`

	MatchServerVersion bool `yaml:"match_server_version"`

	BasicAuth struct {
		rs.BaseField `yaml:"-"`

		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"basic_auth"`

	TLS struct {
		rs.BaseField `yaml:"-"`

		CaCert string `yaml:"ca_cert"`

		InsecureSkipVerify bool `yaml:"insecure_skip_verify"`

		ClientCert string `yaml:"client_cert"`
		ClientKey  string `yaml:"client_key"`
	} `yaml:"tls"`

	Impersonate []kubeContext `yaml:"impersonate"`
}
