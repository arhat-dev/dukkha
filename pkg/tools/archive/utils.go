package archive

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/afero"
)

type entry struct {
	to   string
	from string
	info fs.FileInfo

	link string
}

// collectFiles to be archived
func collectFiles(files []*archiveFileSpec) ([]*entry, error) {
	var (
		candidates []*entry

		ret []*entry
	)

	_fs := afero.NewIOFS(afero.NewOsFs())
	for _, f := range files {
		from := filepath.Clean(f.From)
		actualPaths, err := doublestar.Glob(_fs, from)
		if err != nil {
			actualPaths = []string{from}
		}

		// normalize matched paths
		slashMatches := make([]string, len(actualPaths))
		for i, v := range actualPaths {
			if len(v) == 0 {
				v = "."
			}

			info, err := os.Lstat(v)
			if err != nil {
				return nil, err
			}

			link := ""
			if info.Mode()&fs.ModeSymlink != 0 {
				link, err = os.Readlink(v)
				if err != nil {
					return nil, err
				}
			}

			candidates = append(candidates, &entry{
				info: info,
				from: v,
				link: link,
			})

			slashMatches[i] = filepath.ToSlash(v)
			if info.IsDir() {
				slashMatches[i] += "/"
			}
		}

		size := len(candidates)

		dstIsDir := strings.HasSuffix(f.To, "/") || len(f.To) == 0
		switch len(slashMatches) {
		case 0:
			return nil, fmt.Errorf("no file matches pattern %q", f.From)
		case 1:
			toAdd := candidates[size-1]
			// only one file match
			if dstIsDir {
				toAdd.to = path.Join(f.To, filepath.Base(toAdd.from))
			} else {
				toAdd.to = f.To
			}

			ret = append(ret, toAdd)
			continue
		default:
		}

		// multiple file matches, expect dir for these files
		if !dstIsDir {
			return nil, fmt.Errorf("too many files to %q", f.To)
		}

		prefix := lcpp(slashMatches)
		for i, slashPath := range slashMatches {
			if slashPath == prefix {
				continue
			}

			toAdd := candidates[size-len(slashMatches)+i]
			toAdd.to = path.Join(f.To, strings.TrimPrefix(slashPath, prefix))

			ret = append(ret, toAdd)
		}
	}

	return ret, nil
}

func copyFileContent(w io.Writer, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}

	return nil
}

// longest common path prefix
func lcpp(l []string) string {
	if len(l) == 0 {
		return ""
	}

	min, max := l[0], l[0]
	for _, s := range l[1:] {
		switch {
		case s < min:
			min = s
		case s > max:
			max = s
		}
	}

	lastSlashAt := -1
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			break
		}

		if min[i] == '/' {
			lastSlashAt = i
		}
	}

	if lastSlashAt != -1 {
		return min[:lastSlashAt+1]
	}

	return ""
}
