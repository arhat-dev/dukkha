package conf

import (
	"errors"
	"fmt"
	"io/fs"
	"path"

	"arhat.dev/pkg/log"
	ds "github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

func ReadConfigRecursively(
	rootfs fs.FS,
	configPaths []string,
	ignoreFileNotExist bool,
	visitedPaths *map[string]struct{},
	mergedConfig *Config,
) error {
	for _, target := range configPaths {
		info, err := fs.Stat(rootfs, target)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				if ignoreFileNotExist {
					continue
				}
			}

			return err
		}

		if !info.IsDir() {
			err = readAndMergeConfigFile(
				rootfs, visitedPaths, mergedConfig, target,
			)
			if err != nil {
				return err
			}

			continue
		}

		dirFS, err := fs.Sub(rootfs, target)
		if err != nil {
			return err
		}

		err = fs.WalkDir(dirFS, ".", func(pathInDir string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			switch path.Ext(pathInDir) {
			case ".yaml":
				// leave .yml for customization
			default:
				return nil
			}

			return readAndMergeConfigFile(
				rootfs, visitedPaths, mergedConfig, path.Join(target, pathInDir),
			)
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func readAndMergeConfigFile(
	rootfs fs.FS,
	visitedPaths *map[string]struct{},
	mergedConfig *Config,
	file string,
) error {
	if _, ok := (*visitedPaths)[file]; ok {
		return nil
	}

	(*visitedPaths)[file] = struct{}{}

	configBytes, err := fs.ReadFile(rootfs, file)
	if err != nil {
		return fmt.Errorf("read config file %q: %w", file, err)
	}

	current := NewConfig()
	err = yaml.Unmarshal(configBytes, &current)
	if err != nil {
		return fmt.Errorf("unmarshal config file %q: %w", file, err)
	}

	log.Log.V("config unmarshaled", log.String("file", file), log.Any("config", current))

	err = mergedConfig.Merge(current)
	if err != nil {
		return fmt.Errorf("merge config file %q: %w", file, err)
	}

	for _, inc := range current.Include {
		log.Log.V("working on include entry", log.String("value", inc))

		var toInclude string
		if path.IsAbs(inc) {
			toInclude = inc
		} else {
			toInclude = path.Join(path.Dir(file), inc)
		}

		matches, err2 := ds.Glob(rootfs, toInclude)
		if err2 != nil {
			matches = []string{toInclude}
		}

		err2 = ReadConfigRecursively(
			rootfs, matches, false, visitedPaths, mergedConfig,
		)
		if err2 != nil {
			return fmt.Errorf("loading included config files: %w", err2)
		}
	}

	return err
}
