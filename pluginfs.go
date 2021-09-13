package dukkha

import (
	"embed"
	"io/fs"
	"path"

	"arhat.dev/unionfs"
)

var (
	//go:embed pkg/dukkha/rendering.go
	//go:embed pkg/dukkha/rendering_values.go
	//go:embed pkg/dukkha/tasks.go
	//go:embed pkg/dukkha/tools.go
	//go:embed pkg/dukkha/types.go
	//go:embed pkg/dukkha/context.go
	//go:embed pkg/dukkha/context_exec.go
	//go:embed pkg/dukkha/context_std.go
	//go:embed pkg/dukkha/shells.go
	//go:embed pkg/matrix
	//go:embed pkg/renderer/*.go
	//go:embed pkg/sliceutils/*.go
	//go:embed pkg/tools/action.go
	//go:embed pkg/tools/hook.go
	//go:embed pkg/tools/task.go
	//go:embed pkg/tools/tool.go
	//go:embed pkg/tools/utils.go
	//go:embed pkg/tools/*_pseudo.go
	//go:embed pkg/constant/*.go
	pkg_fs embed.FS

	//go:embed third_party/arhat.dev/rs
	vendor_rs_fs embed.FS

	//go:embed third_party/gopkg.in/yaml.v3
	vendor_yaml_fs embed.FS

	//go:embed third_party/github.com/evanphx/json-patch/v5
	vendor_jsonpatch_fs embed.FS

	//go:embed third_party/github.com/pkg/errors
	vendor_pkgerrors_fs embed.FS

	//go:embed third_party/mvdan.cc/sh/v3
	vendor_mvdansh_fs embed.FS

	//go:embed third_party/arhat.dev/pkg
	vendor_arhatpkg_fs embed.FS
)

func NewPluginFS(goPath, pkg string) fs.FS {
	basePath := path.Join(goPath, "src", pkg, "vendor")

	pfs := unionfs.New()

	pfs.Map(path.Join(basePath, "arhat.dev/dukkha/pkg"), "pkg", pkg_fs)

	for vendorPkg, vendorPkgFS := range map[string]embed.FS{
		"arhat.dev/rs":                     vendor_rs_fs,
		"gopkg.in/yaml.v3":                 vendor_yaml_fs,
		"github.com/evanphx/json-patch/v5": vendor_jsonpatch_fs,
		"github.com/pkg/errors":            vendor_pkgerrors_fs,
		"mvdan.cc/sh/v3":                   vendor_mvdansh_fs,
		// "github.com/huandu/xstrings":       vendor_xstring_fs,
		// "github.com/muesli/termenv":        vendor_termenv_fs,
		// "github.com/lucasb-eyer/go-colorful": vendor_colorful_fs,
		// "golang.org/x/sys":                   vendor_sys_fs,
		// "golang.org/x/sync":                  vendor_sync_fs,
		"arhat.dev/pkg": vendor_arhatpkg_fs,
		// "github.com/itchyny/gojq":            vendor_gojq_fs,
		// "github.com/itchyny/timefmt-go": vendor_timefmt_fs,
	} {
		pfs.Map(
			path.Join(basePath, vendorPkg),
			path.Join("third_party", vendorPkg),
			vendorPkgFS,
		)
	}

	return pfs
}
