package archive

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"arhat.dev/pkg/fshelper"
)

type entry struct {
	to   string
	from string
	link string // symlink target
	info fs.FileInfo
}

const maxSymlinkDepth = 256

// collectFiles to be archived
func collectFiles(
	ofs *fshelper.OSFS,
	specs []*fileFromToSpec,
) (ret []*entry, err error) {
	var (
		// swap first serves as the storage of all matched path info
		swap []*entry

		srcPaths []string

		// key: in archive path, dir entry with trailing slash added
		// value: index into ret
		inArchiveFiles = make(map[string]int)
	)

	for _, spec := range specs {
		// step 1: get all matched filesystem paths
		// NOTE: do not do filepath.ToSlash(spec.From) to require slash style pattern
		srcPaths, err = ofs.Glob(spec.From)
		if err != nil {
			err = nil
			srcPaths = []string{spec.From}
		}

		// step 2: normalize matched paths, store necessary path info
		// 1. enforce unix style path (`/`)
		// 2. add trailing slash to dir entries
		slashMatches := make([]string, len(srcPaths))
		for i, sp := range srcPaths {
			if len(sp) == 0 {
				sp = "."
			}

			var (
				info fs.FileInfo
				link string

				symlinkDepth int
			)

			info, err = ofs.Lstat(sp)
			if err != nil {
				return
			}

		checkSymlink:
			symlinkDepth++
			if symlinkDepth > maxSymlinkDepth {
				err = fmt.Errorf("too many follow symlink operations for one file (> 256)")
				return
			}

			if info.Mode()&fs.ModeSymlink != 0 {
				link, err = ofs.Readlink(sp)
				if err != nil {
					return
				}

				if spec.FollowSymlink {
					sp = link
					link = ""
					info, err = ofs.Lstat(sp)
					if err != nil {
						return
					}

					goto checkSymlink
				}
			}

			swap = append(swap, &entry{
				info: info,
				from: sp,
				link: link,
			})

			slashMatches[i] = filepath.ToSlash(sp)
			if info.IsDir() {
				slashMatches[i] += "/"
			}
		}

		szSwap := len(swap)

		nSlashMatches := len(slashMatches)
		dstIsDir := strings.HasSuffix(spec.To, "/") || len(spec.To) == 0
		switch nSlashMatches {
		case 0:
			err = fmt.Errorf("no file matching pattern %q", spec.From)
			return
		case 1:
			toAdd := swap[szSwap-1]
			// only one file match
			if dstIsDir {
				toAdd.to = path.Join(spec.To, filepath.Base(toAdd.from))
			} else {
				toAdd.to = spec.To
			}

			inArchiveFiles[toAdd.to] = len(ret)
			ret = append(ret, toAdd)
			continue
		default:
		}

		// multiple file matches, expect dir for these files
		if !dstIsDir {
			err = fmt.Errorf("multiple files selected, expecting a dir, got %q, missing trailing slash?", spec.To)
			return
		}

		prefix := lcpp(slashMatches)
		for i, sp := range slashMatches {
			if sp == prefix {
				continue
			}

			toAdd := swap[szSwap-nSlashMatches+i]
			toAdd.to = path.Join(spec.To, strings.TrimPrefix(sp, prefix))

			inArchiveFiles[toAdd.to] = len(ret)
			ret = append(ret, toAdd)
		}
	}

	// add missing directories

	// good is the set of in archive paths having a viable parent
	// key: in archive path
	good := make(map[string]struct{})
	nLastKnownGood := -1
	for nLastKnownGood != len(good) {
		nLastKnownGood = len(good)
		for k, idx := range inArchiveFiles {
			dir := path.Dir(k)
			if dir == "." || dir == "/" {
				good[k] = struct{}{}
				continue
			}

			// TODO: why this exists?
			_, ok := inArchiveFiles[dir]
			if ok {
				good[k] = struct{}{}
				continue
			}

			_, ok = inArchiveFiles[dir+"/"]
			if ok {
				good[k] = struct{}{}
				continue
			}

			// no parent dir, add a fake one based on
			// actual parent of the file
			from := filepath.Dir(ret[idx].from)
			info, err := ofs.Lstat(from)
			if err != nil {
				return nil, err
			}

			ent := &entry{
				from: from,
				info: info,
				to:   dir,
				link: "",
			}

			inArchiveFiles[dir] = len(ret)
			ret = append(ret, ent)
		}
	}

	sort.Slice(ret, func(i, j int) bool { return ret[i].to < ret[j].to })
	return ret, nil
}

func copyFileContent(ofs *fshelper.OSFS, w io.Writer, file string) error {
	f, err := ofs.Open(file)
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
