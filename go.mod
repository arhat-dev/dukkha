module arhat.dev/dukkha

go 1.18

require (
	arhat.dev/pkg v0.10.2-0.20220730124454-80bcc6fbd51f
	arhat.dev/rs v0.10.1-0.20220730165847-4b9ceeed1310
	arhat.dev/tlang v0.0.0-20220721161405-faac9d553010
	github.com/Masterminds/goutils v1.1.1
	github.com/aoldershaw/ansi v0.0.0-20210128170437-8c5426635e02
	github.com/bmatcuk/doublestar/v4 v4.2.0
	github.com/d5/tengo/v2 v2.12.1
	github.com/die-net/lrucache v0.0.0-20220628165024-20a71bc65bf1
	github.com/dsnet/compress v0.0.1
	github.com/google/uuid v1.3.0
	github.com/gosimple/slug v1.12.0
	github.com/h2non/filetype v1.1.3
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/huandu/xstrings v1.3.2
	github.com/itchyny/gojq v0.12.8
	github.com/klauspost/compress v1.15.9
	github.com/minio/minio-go/v7 v7.0.31
	github.com/mitchellh/copystructure v1.2.0
	github.com/muesli/termenv v0.12.0
	github.com/nwaples/rardecode v1.1.3
	github.com/open2b/scriggo v0.55.0
	github.com/pierrec/lz4/v4 v4.1.15
	github.com/spf13/cobra v1.5.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.0
	github.com/ulikunitz/xz v0.5.10
	github.com/weaveworks/schemer v0.0.0-20210802122110-338b258ad2ca
	go.uber.org/multierr v1.8.0
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f
	golang.org/x/term v0.0.0-20220722155259-a9ba230a4035
	golang.org/x/tools v0.1.11
	gopkg.in/yaml.v3 v3.0.1
	mvdan.cc/sh/v3 v3.5.1
)

replace (
	// branch master
	github.com/d5/tengo/v2 => github.com/arhat-dev/tengo/v2 v2.12.1-0.20220719235156-be261837907d
	// branch master
	github.com/weaveworks/schemer => github.com/arhat-dev/schemer v0.0.0-20220730131244-3259af3d11b8
	// branch `dukkha`
	mvdan.cc/sh/v3 => github.com/arhat-dev/sh/v3 v3.5.0-0.dev.0.20220723194407-b237320fa712
)

require (
	arhat.dev/pty v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/itchyny/timefmt-go v0.1.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/minio/md5-simd v1.1.0 // indirect
	github.com/minio/sha256-simd v0.1.1 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rs/xid v1.2.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/smartystreets/assertions v1.13.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
)
