//go:build !noself

package dukkha_internal

import (
	"io"
	"strings"

	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/dukkha"

	_ "unsafe" // for go:linkname
)

//go:linkname newDukkhaCmd arhat.dev/dukkha/pkg/cmd.NewRootCmd
func newDukkhaCmd(prevCtx dukkha.Context) *cobra.Command

// RunSelf runs dukkha itself without spawning a new process
func RunSelf(
	ctx dukkha.Context,
	stdin io.Reader,
	stdout, stderr io.Writer,
	args ...string,
) (err error) {
	ctx = ctx.DeriveNew()
	ctx.SetStdIO(stdin, stdout, stderr)

	cmd := newDukkhaCmd(ctx)
	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// set cmd.args (unexported field of cobra.Command) for cmd.Execute
	cmd.SetArgs(args)

	// only skip config resolving when --config/-c not specified
	resolveConfig := false
	for _, v := range args {
		switch {
		case strings.HasPrefix(v, "--config"), strings.HasPrefix(v, "-c"):
			// do not skip config resovling
			resolveConfig = true
		}

		if resolveConfig {
			break
		}
	}

	if !resolveConfig {
		cmd.PersistentPreRunE = nil
	}

	return cmd.ExecuteContext(ctx)
}
