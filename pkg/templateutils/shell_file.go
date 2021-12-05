package templateutils

import (
	"context"
	"io"
	"io/fs"
	"os"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/iohelper"
	"mvdan.cc/sh/v3/interp"
)

func fileOpenHandler(
	ctx context.Context,
	path string,
	flag int,
	perm fs.FileMode,
) (io.ReadWriteCloser, error) {
	const devNullPath = "/dev/null"

	if path == devNullPath {
		return iohelper.NewDevNull(), nil
	}

	hc := interp.HandlerCtx(ctx)

	ofs := fshelper.NewOSFS(false, func() (string, error) {
		return hc.Dir, nil
	})

	f, err := ofs.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}

	return f.(*os.File), nil
}
