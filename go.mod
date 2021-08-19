module arhat.dev/dukkha

go 1.16

require (
	arhat.dev/pkg v0.6.0
	arhat.dev/rs v0.1.1
	github.com/Masterminds/goutils v1.1.1
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/aoldershaw/ansi v0.0.0-20210128170437-8c5426635e02
	github.com/die-net/lrucache v0.0.0-20210801000212-e34e67316dc5
	github.com/google/uuid v1.3.0
	github.com/gosimple/slug v1.10.0
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/huandu/xstrings v1.3.2
	github.com/itchyny/gojq v0.12.4
	github.com/muesli/termenv v0.9.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/multierr v1.7.0
	golang.org/x/crypto v0.0.0-20210813211128-0a44fdfbc16e
	golang.org/x/sys v0.0.0-20210816074244-15123e1e1f71
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	mvdan.cc/sh/v3 v3.3.1
)

replace (
	arhat.dev/rs => arhat.dev/rs v0.1.1
	github.com/creack/pty => github.com/jeffreystoke/pty v1.1.12-0.20210531091229-b834701fbcc6
)
