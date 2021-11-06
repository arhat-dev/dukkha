module arhat.dev/dukkha

go 1.16

require (
	arhat.dev/pkg v0.7.3
	arhat.dev/rs v0.7.1-0.20211106175455-04fd26f5a52a
	github.com/Masterminds/goutils v1.1.1
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/aoldershaw/ansi v0.0.0-20210128170437-8c5426635e02
	github.com/die-net/lrucache v0.0.0-20210908122246-903d43d14082
	github.com/google/uuid v1.3.0
	github.com/gosimple/slug v1.11.0
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/huandu/xstrings v1.3.2
	github.com/itchyny/gojq v0.12.5
	github.com/minio/minio-go/v7 v7.0.15
	github.com/muesli/termenv v0.9.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/weaveworks/schemer v0.0.0-20210802122110-338b258ad2ca
	go.uber.org/multierr v1.7.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/net v0.0.0-20211104170005-ce137452f963
	golang.org/x/sys v0.0.0-20211103235746-7861aae1554b
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	golang.org/x/tools v0.1.7
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	mvdan.cc/sh/v3 v3.4.0
)

replace (
	github.com/creack/pty => github.com/donorp/pty v1.1.12-0.20211004111936-294eccab62ed
	github.com/weaveworks/schemer => github.com/arhat-dev/schemer v0.0.0-20211102163138-8bc12e169191
)
