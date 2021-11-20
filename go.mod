module arhat.dev/dukkha

go 1.16

require (
	arhat.dev/pkg v0.8.1-0.20211116085804-2e58eae7aa02
	arhat.dev/rs v0.8.4-0.20211113095115-c107aceb0d05
	github.com/Masterminds/goutils v1.1.1
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/aoldershaw/ansi v0.0.0-20210128170437-8c5426635e02
	github.com/bmatcuk/doublestar/v4 v4.0.2
	github.com/die-net/lrucache v0.0.0-20210908122246-903d43d14082
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5
	github.com/google/uuid v1.3.0
	github.com/gosimple/slug v1.11.2
	github.com/h2non/filetype v1.1.2-0.20210917125640-7fafb18134ff
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/huandu/xstrings v1.3.2
	github.com/itchyny/gojq v0.12.5
	github.com/klauspost/compress v1.13.6
	github.com/minio/minio-go/v7 v7.0.15
	github.com/muesli/termenv v0.9.0
	github.com/nwaples/rardecode v1.1.2
	github.com/pierrec/lz4/v4 v4.1.10
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/ulikunitz/xz v0.5.10
	github.com/weaveworks/schemer v0.0.0-20210802122110-338b258ad2ca
	go.uber.org/multierr v1.7.0
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa
	golang.org/x/net v0.0.0-20211111160137-58aab5ef257a
	golang.org/x/sys v0.0.0-20211110154304-99a53858aa08
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
