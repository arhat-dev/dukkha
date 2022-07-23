package conf

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"

	ds "github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
)

// Read config in rootfs recursively
func Read(
	rc dukkha.ConfigResolvingContext,
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
			err = readAndMergeConfigFile(rc,
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

			return readAndMergeConfigFile(rc,
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
	rc dukkha.ConfigResolvingContext,
	rootfs fs.FS,
	visitedPaths *map[string]struct{},
	mergedConfig *Config,
	file string,
) error {
	if _, ok := (*visitedPaths)[file]; ok {
		return nil
	}

	(*visitedPaths)[file] = struct{}{}

	r, err := rootfs.Open(file)
	if err != nil {
		return fmt.Errorf("read config file %q: %w", file, err)
	}

	include, err := loadConfig(rc, r, mergedConfig)
	_ = r.Close()
	if err != nil {
		return err
	}

	return handleInclude(rc, rootfs, visitedPaths, mergedConfig, file, include)
}

// loadConfig unmarshal all yaml docs in r as Config, add resolved renderers into rc
// then merge freshly unmarshaled Config into mergedConfig
func loadConfig(
	rc dukkha.ConfigResolvingContext,
	r io.Reader,
	mergedConfig *Config,
) ([]*IncludeEntry, error) {
	var ret []*IncludeEntry

	dec := yaml.NewDecoder(r)
	for {
		current := NewConfig()

		err := dec.Decode(current)
		if err != nil {
			if err == io.EOF {
				return ret, nil
			}

			return nil, fmt.Errorf("unmarshal config: %w", err)
		}

		err = current.resolveRenderers(rc)
		if err != nil {
			return nil, fmt.Errorf("resolve renderers: %w", err)
		}

		err = current.resolveShells(rc)
		if err != nil {
			return nil, fmt.Errorf("resolve shells: %w", err)
		}

		err = current.ResolveFields(rc, -1, "include")
		if err != nil {
			return nil, fmt.Errorf("resolve include entries: %w", err)
		}

		ret = append(ret, current.Include...)

		err = mergedConfig.Merge(current)
		if err != nil {
			return nil, fmt.Errorf("merge config: %w", err)
		}
	}
}

func handleInclude(
	rc dukkha.ConfigResolvingContext,
	rootfs fs.FS,
	visitedPaths *map[string]struct{},
	mergedConfig *Config,
	currentFile string,
	include []*IncludeEntry,
) error {
	for _, inc := range include {
		switch {
		case len(inc.Path) != 0:
			toInclude := inc.Path
			if !path.IsAbs(toInclude) {
				// TODO: whether relative current file or DUKKHA_WORKDIR
				toInclude = path.Join(path.Dir(currentFile), toInclude)
			}

			matches, err2 := ds.Glob(rootfs, toInclude)
			if err2 != nil {
				matches = []string{toInclude}
			}

			err2 = Read(rc,
				rootfs, matches, false, visitedPaths, mergedConfig,
			)

			if err2 != nil {
				return fmt.Errorf("loading included config files: %w", err2)
			}
		case len(inc.Text) != 0:
			embedInclude, err := loadConfig(rc, strings.NewReader(inc.Text), mergedConfig)
			if err != nil {
				return err
			}

			err = handleInclude(rc, rootfs, visitedPaths, mergedConfig, currentFile, embedInclude)
			if err != nil {
				return err
			}
		default:
			continue
		}
	}

	return nil
}
