package render

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
)

func renderYamlFile(
	rc dukkha.Context,
	srcPath string,
	destPath *string,
	opts *ResolvedOptions,
	srcPerm map[string]os.FileMode,
) error {
	sInfo, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to check src dir %q: %w", srcPath, err)
	}

	srcPerm[srcPath] = sInfo.Mode().Perm()

	if sInfo.IsDir() {
		entries, err2 := os.ReadDir(srcPath)
		if err2 != nil {
			return fmt.Errorf("failed to check files in src dir %q: %w", srcPath, err2)
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

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	// srcPath exitsts, ensure destination if any

	if destPath != nil {
		// prepare destination parent dir if not exists

		err = ensureDestDir(srcPath, *destPath, srcPerm)
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

func ensureDestDir(srcPath, destPath string, srcPerm map[string]os.FileMode) error {
	srcPath = filepath.Dir(srcPath)
	destPath = filepath.Dir(destPath)
	var doMkdir []func() error

	for {
		_, err := os.Stat(destPath)
		if err == nil {
			// already exists, do nothing
			break
		}

		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to check dest dir %q: %w", destPath, err)
		}

		perm, ok := srcPerm[srcPath]
		if !ok {
			// checking parent dir of user priveded src dir
			info, err2 := os.Stat(srcPath)
			if err2 != nil {
				return fmt.Errorf("failed to check src parent dir %q: %w", srcPath, err2)
			}

			perm = info.Mode().Perm()
		}

		// copy value, do not reference srcPath and destDir directly
		targetDir := destPath
		src := srcPath
		doMkdir = append(doMkdir, func() error {
			err = os.Mkdir(targetDir, perm)
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
