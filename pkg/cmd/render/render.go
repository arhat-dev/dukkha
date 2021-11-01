package render

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"arhat.dev/rs"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
)

func renderYamlFileOrDir(
	rc dukkha.Context,
	srcPath string,
	destPath *string,
	outputFormat string,
	createEnc encoderCreateFunc,
	recursive bool,
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
				err3 = renderYamlFileOrDir(
					rc, newSrc, newDest, outputFormat,
					createEnc, recursive, srcPerm,
				)
				if err3 != nil {
					return fmt.Errorf("failed to render file %q: %w", newSrc, err3)
				}
			default:
				if ent.IsDir() && recursive {
					err3 = renderYamlFileOrDir(
						rc, newSrc, newDest, outputFormat,
						createEnc, recursive, srcPerm,
					)
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

	dec := yaml.NewDecoder(srcFile)
	ret, err := parseYaml(dec, func() interface{} { return new(rs.AnyObject) })

	// always write parsed yaml docs, regardless of errors

	var enc encoder
	if len(ret) != 0 {
		if destPath == nil {
			enc, err = createEnc(os.Stdout)
			if err != nil {
				return err
			}
		} else {
			// prepare destination parent dir if not exists

			err = ensureDestDir(srcPath, *destPath, srcPerm)
			if err != nil {
				return err
			}

			dest := *destPath
			switch outputFormat {
			case "json":
				// change extension name
				dest = strings.TrimSuffix(dest, filepath.Ext(srcPath)) + ".json"
			case "yaml":
				// do nothing since source file is yaml as well
				// no matter .yml or .yaml
			}
			var destFile *os.File
			destFile, err = os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcPerm[srcPath])
			if err != nil {
				return fmt.Errorf("failed to open output file %q: %w", dest, err)
			}
			defer func() { _ = destFile.Close() }()

			enc, err = createEnc(destFile)
			if err != nil {
				return err
			}
		}
	}

	err = os.Chdir(filepath.Dir(srcPath))
	if err != nil {
		return fmt.Errorf("failed to change working dir to source parent: %w", err)
	}
	defer func() {
		err = os.Chdir(rc.WorkingDir())
		if err != nil {
			panic(fmt.Errorf(
				"failed to go back to dukkha working dir: %w", err,
			))
		}
	}()

	for _, doc := range ret {
		obj := doc.(*rs.AnyObject)

		err2 := obj.ResolveFields(rc, -1)
		if err2 != nil {
			return multierr.Append(err, err2)
		}

		err2 = enc.Encode(obj)
		if err2 != nil {
			return multierr.Append(err, err2)
		}
	}

	return err
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

func parseYaml(dec *yaml.Decoder, createOutObj func() interface{}) ([]interface{}, error) {
	var ret []interface{}
	for {
		obj := createOutObj()
		err := dec.Decode(obj)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return ret, nil
			}

			return ret, err
		}

		if obj == nil {
			continue
		}

		ret = append(ret, obj)
	}
}
