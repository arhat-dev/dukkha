package archive

import (
	"archive/tar"
	"io"
	"strings"

	"arhat.dev/pkg/fshelper"
)

func createTar(
	ofs *fshelper.OSFS,
	w io.Writer,
	files []*entry,
	enableCompression bool,
	compressionMethod string,
	compressionLevel string,
) (err error) {
	var (
		tw  *tar.Writer
		hdr *tar.Header
	)

	if enableCompression {
		var tarball io.WriteCloser
		tarball, err = createCompressionStream(w, compressionMethod, compressionLevel)
		if err != nil {
			return
		}

		tw = tar.NewWriter(tarball)
		defer func() {
			_ = tw.Close()
			_ = tarball.Close()
		}()
	} else {
		tw = tar.NewWriter(w)
		defer func() { _ = tw.Close() }()
	}

	for _, f := range files {
		hdr, err = tar.FileInfoHeader(f.info, f.link)
		if err != nil {
			return
		}

		hdr.Format = tar.FormatPAX
		hdr.Name = f.to

		mode := f.info.Mode()
		if mode.IsDir() && !strings.HasSuffix(hdr.Name, "/") {
			hdr.Name += "/"
		}

		err = tw.WriteHeader(hdr)
		if err != nil {
			return
		}

		if mode.IsRegular() {
			err = copyFileContent(ofs, tw, f.from)
			if err != nil {
				return
			}
		}

		err = tw.Flush()
		if err != nil {
			return
		}
	}

	return
}
