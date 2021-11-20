package archive

import (
	"archive/tar"
	"io"
	"strings"
)

func createTar(w io.Writer, files []*entry) error {
	tw := tar.NewWriter(w)
	defer func() { _ = tw.Close() }()

	for _, f := range files {
		hdr, err := tar.FileInfoHeader(f.info, f.link)
		if err != nil {
			return err
		}

		hdr.Format = tar.FormatPAX
		hdr.Name = f.to
		if f.info.IsDir() && !strings.HasSuffix(f.to, "/") {
			f.to += "/"
		}

		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}

		if f.info.Mode().IsRegular() {
			err = copyFileContent(tw, f.from)
			if err != nil {
				return err
			}
		}

		err = tw.Flush()
		if err != nil {
			return err
		}
	}

	return nil
}
