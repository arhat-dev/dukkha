package af

import (
	"archive/tar"
	"fmt"
	"io"
	"path"
)

// when target is empty or `.`, return reader of first regular file
// in the tarball
func untar(r io.ReadSeeker, target string) (io.Reader, error) {
	restore, err := prepareSeekRestore(r)
	if err != nil {
		return nil, err
	}

	rd := tar.NewReader(r)
	for {
		hdr, err := rd.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("untar: %w", err)
		}

		name := path.Clean(hdr.Name)
		switch {
		case target == ".",
			len(target) == 0:
			// empty target, grab first regular file
			switch hdr.Typeflag {
			case tar.TypeReg, tar.TypeRegA:
				return rd, nil
			default:
				continue
			}
		case name != target:
			continue
		}

		switch hdr.Typeflag {
		case tar.TypeReg, tar.TypeRegA:
			return rd, nil
		case tar.TypeLink, tar.TypeSymlink:
			err = restore()
			if err != nil {
				return nil, err
			}

			// only allow internal link
			return untar(r, cleanLink(name, hdr.Linkname))
		default:
			return nil, fmt.Errorf("untar: unsupported non regular file %q", name)
		}
	}

	return nil, fmt.Errorf("untar: file %q not found in archive", target)
}
