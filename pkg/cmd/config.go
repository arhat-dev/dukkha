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

type configLoaderFunc func(
	visitedPaths *map[string]struct{},
	mergedConfig interface{},
	file string,
	loader configLoaderFunc,
) error

func readAndMergeConfigFile(
	visitedPaths *map[string]struct{},
	mergedConfig interface{},
	file string,
	loader configLoaderFunc,
) error {
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

	err = mergedConfig.(*conf.Config).Merge(current)
	if err != nil {
		return fmt.Errorf("failed to merge config file %q: %w", file, err)
	}

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

		err2 = readConfigRecursively(matches, false, visitedPaths, mergedConfig, loader)
		if err2 != nil {
			return fmt.Errorf("failed to load included config files: %w", err2)
		}
	}

	return err
}

func readAndMergePluginConfigFile(
	visitedPaths *map[string]struct{},
	mergedConfig interface{},
	file string,
	loader configLoaderFunc,
) error {
	if _, ok := (*visitedPaths)[file]; ok {
		return nil
	}

	(*visitedPaths)[file] = struct{}{}

	configBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read plugin config file %q: %w", file, err)
	}

	current := conf.NewPluginConfig()
	err = yaml.Unmarshal(configBytes, &current)
	if err != nil {
		return fmt.Errorf("failed to unmarshal plugin config file %q: %w", file, err)
	}

	log.Log.V("config unmarshaled", log.String("file", file), log.Any("config", current))

	err = mergedConfig.(*conf.PluginConfig).Merge(current)
	if err != nil {
		return fmt.Errorf("failed to merge plugin config %q: %w", file, err)
	}

	return err
}

func readConfigRecursively(
	configPaths []string,
	ignoreFileNotExist bool,
	visitedPaths *map[string]struct{},
	mergedConfig interface{},
	loader configLoaderFunc,
) error {
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
			err = loader(visitedPaths, mergedConfig, path, loader)
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

			return loader(
				visitedPaths,
				mergedConfig,
				filepath.Join(path, pathInDir),
				loader,
			)
		})

		if err != nil {
			return err
		}
	}

	return nil
}
