package render

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/pkg/fshelper"
)

func renderYamlFile(
	rc dukkha.Context,
	srcPath string,
	destPath *string,
	opts *ResolvedOptions,
	srcPerm map[string]fs.FileMode,
) error {
	sInfo, err := rc.FS().Stat(srcPath)
	if err != nil {
		return fmt.Errorf("invalid file source: %w", err)
	}

	srcPerm[srcPath] = sInfo.Mode().Perm()

	if sInfo.IsDir() {
		entries, err2 := rc.FS().ReadDir(srcPath)
		if err2 != nil {
			return fmt.Errorf("unable to check files in src dir %q: %w", srcPath, err2)
		}

		for _, ent := range entries {
			name := ent.Name()
			newSrc := filepath.Join(srcPath, name)

			var newDest *string
			if destPath != nil {
				newDestPath := filepath.Join(*destPath, name)
				newDest = &newDestPath
			}

			newSrcInfo, err3 := ent.Info()
			if err3 != nil {
				return fmt.Errorf("failed to check dir entry %q: %w", newSrc, err3)
			}

			srcPerm[newSrc] = newSrcInfo.Mode().Perm()

			switch filepath.Ext(name) {
			case ".yml", ".yaml":
				err3 = renderYamlFile(rc, newSrc, newDest, opts, srcPerm)
				if err3 != nil {
					return fmt.Errorf("failed to render file %q: %w", newSrc, err3)
				}
			default:
				if ent.IsDir() && opts.recursive {
					err3 = renderYamlFile(rc, newSrc, newDest, opts, srcPerm)
					if err3 != nil {
						return fmt.Errorf("failed to render dir %q: %w", newSrc, err3)
					}
				}
			}
		}

		return nil
	}

	// srcPath should be a yaml file

	srcFile, err := rc.FS().Open(srcPath)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	// srcPath exitsts, ensure destination if any

	if destPath != nil {
		// prepare destination parent dir if not exists

		err = ensureDestDir(rc.FS(), srcPath, *destPath, srcPerm)
		if err != nil {
			return err
		}

		dest := *destPath
		switch opts.outputFormat {
		case "json":
			// change extension name
			dest = strings.TrimSuffix(dest, filepath.Ext(srcPath)) + ".json"
		case "yaml":
			// do nothing since source file is yaml as well
			// no matter .yml or .yaml
		}

		destPath = &dest
	}

	return renderYamlReader(rc, srcFile, destPath, srcPerm[srcPath], opts)
}

func ensureDestDir(ofs *fshelper.OSFS, srcPath, destPath string, srcPerm map[string]fs.FileMode) error {
	srcPath = filepath.Dir(srcPath)
	destPath = filepath.Dir(destPath)
	var doMkdir []func() error

	for {
		_, err := ofs.Stat(destPath)
		if err == nil {
			// already exists, do nothing
			break
		}

		if !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("failed to check dest dir %q: %w", destPath, err)
		}

		perm, ok := srcPerm[srcPath]
		if !ok {
			// checking parent dir of user priveded src dir
			info, err2 := ofs.Stat(srcPath)
			if err2 != nil {
				return fmt.Errorf("failed to check src parent dir %q: %w", srcPath, err2)
			}

			perm = info.Mode().Perm()
		}

		// copy value, do not reference srcPath and destDir directly
		targetDir := destPath
		src := srcPath
		doMkdir = append(doMkdir, func() error {
			err = ofs.Mkdir(targetDir, perm)
			if err != nil {
				return fmt.Errorf(
					"failed to create dest dir %q for src dir %q: %w",
					targetDir, src, err,
				)
			}

			return nil
		})

		srcPath = filepath.Dir(srcPath)
		destPath = filepath.Dir(destPath)
	}

	for i := len(doMkdir) - 1; i >= 0; i-- {
		err := doMkdir[i]()
		if err != nil {
			return err
		}
	}

	return nil
}
