package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"arhat.dev/pkg/log"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/conf"
)

func readConfigRecursively(
	configPaths []string,
	ignoreFileNotExist bool,
	visitedPaths *map[string]struct{},
	mergedConfig *conf.Config,
) error {
	readAndMergeConfigFile := func(file string) error {
		if _, ok := (*visitedPaths)[file]; ok {
			return nil
		}

		(*visitedPaths)[file] = struct{}{}

		configBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read config file %q: %w", file, err)
		}

		current := conf.NewConfig()
		err = yaml.Unmarshal(configBytes, &current)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config file %q: %w", file, err)
		}

		log.Log.V("config unmarshaled", log.String("file", file), log.Any("config", current))

		mergedConfig.Merge(current)

		for _, inc := range current.Include {
			log.Log.V("working on include entry", log.String("value", inc))

			var toInclude string
			if filepath.IsAbs(inc) {
				toInclude = inc
			} else {
				toInclude = filepath.Join(filepath.Dir(file), inc)
			}

			matches, err2 := filepath.Glob(toInclude)
			if err2 != nil {
				matches = []string{toInclude}
			}

			err2 = readConfigRecursively(matches, false, visitedPaths, mergedConfig)
			if err2 != nil {
				return fmt.Errorf("failed to load included config files: %w", err)
			}
		}

		return err
	}

	for _, path := range configPaths {
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				if ignoreFileNotExist {
					continue
				}
			}

			return err
		}

		if !info.IsDir() {
			err = readAndMergeConfigFile(path)
			if err != nil {
				return err
			}

			continue
		}

		err = fs.WalkDir(os.DirFS(path), ".", func(pathInDir string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			switch filepath.Ext(pathInDir) {
			case ".yaml":
				// leave .yml for customization
			default:
				return nil
			}

			return readAndMergeConfigFile(filepath.Join(path, pathInDir))
		})

		if err != nil {
			return err
		}
	}

	return nil
}
