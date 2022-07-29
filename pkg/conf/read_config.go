package conf

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
)

type ReadFlag uint32

const (
	ReadFlag_Renderer ReadFlag = 1 << iota
	ReadFlag_Global
	ReadFlag_Tool
	ReadFlag_Task
)

const (
	ReadFlag_Full = ReadFlag_Renderer | ReadFlag_Global | ReadFlag_Tool | ReadFlag_Task
)

type ReadSpec struct {
	Flags        ReadFlag
	ConfFS       fs.FS
	VisitedPaths *map[string]struct{}
	MergedConfig *Config

	lock sync.Mutex
}

// Read config recursively
//
// configPaths are user selected paths (both cli flag --config and yaml include) or defaults,
// if a path entry is a dir, files in that dir with `.yaml` ext are processed in lexical
// order, all matched files are processed in provided order
func Read(
	rc dukkha.ConfigResolvingContext,
	spec *ReadSpec,
	sg *SyncGroup,
	configPaths []string,
	ignoreFileNotExist bool,
) error {
	for i := range configPaths {
		startPath := configPaths[i]

		info, err := fs.Stat(spec.ConfFS, startPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				if ignoreFileNotExist {
					continue
				}
			}

			return err
		}

		switch {
		case info.Mode().IsRegular():
			if !markVisited(spec, startPath) {
				continue
			}

			go readAndMergeConfigFile(rc, spec, sg, startPath, sg.NewJob())

			continue
		case info.IsDir():
			dirFS, err := fs.Sub(spec.ConfFS, startPath)
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

				file := path.Join(startPath, pathInDir)
				if !markVisited(spec, file) {
					return nil
				}

				go readAndMergeConfigFile(rc, spec, sg, file, sg.NewJob())

				return nil
			})

			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid config path %q: not a dir or file", startPath)
		}
	}

	sg.Wait()

	return sg.Err()
}

// return true when the file was not visited
func markVisited(spec *ReadSpec, file string) bool {
	// here we can consider sync.Mutex as spin lock
	spec.lock.Lock()
	if _, ok := (*spec.VisitedPaths)[file]; ok {
		spec.lock.Unlock()
		return false
	}

	(*spec.VisitedPaths)[file] = struct{}{}
	spec.lock.Unlock()
	return true
}

// readAndMergeConfigFile
func readAndMergeConfigFile(
	rc dukkha.ConfigResolvingContext,
	spec *ReadSpec,
	sg *SyncGroup,
	file string,
	j Job,
) {
	r, err := spec.ConfFS.Open(file)
	if err != nil {
		if !sg.Lock(j) {
			sg.Done()
			return
		}

		sg.Cancel(fmt.Errorf("read config file %q: %w", file, err))
		sg.Done()
		return
	}

	loadConfig(rc, spec, sg, r, file, j)
	return
}

// loadConfig unmarshals all yaml docs in r as Config(s), add resolved renderers into rc
// then merge freshly unmarshaled Config into mergedConfig
func loadConfig(
	rc dukkha.ConfigResolvingContext,
	spec *ReadSpec,
	sg *SyncGroup,
	r io.ReadCloser,
	filename string,
	j Job,
) {
	var (
		configs []*Config
		err     error
	)

	defer sg.Done()

	// decode yaml config in parallel
	dec := yaml.NewDecoder(r)
	for {
		current := NewConfig()

		err = dec.Decode(current)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}

			err = fmt.Errorf("decode yaml config: %w", err)
			break
		}

		configs = append(configs, current)
	}
	_ = r.Close()

	// handle config resolving in sequence

	if !sg.Lock(j) {
		return
	}

	if err != nil {
		sg.Cancel(err)
		return
	}

	var (
		includes []*IncludeEntry
		flags    = spec.Flags
	)

	for i, cfg := range configs {
		if flags&ReadFlag_Renderer != 0 {
			err = cfg.resolveRenderers(rc)
			if err != nil {
				sg.Cancel(fmt.Errorf("%s #%d: resolve renderers: %w", filename, i, err))
				return
			}
		}

		err = cfg.resolveShells(rc)
		if err != nil {
			sg.Cancel(fmt.Errorf("%s #%d: resolve shells: %w", filename, i, err))
			return
		}

		err = cfg.ResolveFields(rc, -1, "include")
		if err != nil {
			sg.Cancel(fmt.Errorf("%s #%d: resolve include entries: %w", filename, i, err))
			return
		}

		includes = append(includes, cfg.Include...)

		spec.lock.Lock()
		err = spec.MergedConfig.Merge(cfg)
		if err != nil {
			spec.lock.Unlock()
			sg.Cancel(fmt.Errorf("%s #%d: merge config: %w", filename, i, err))
			return
		}

		spec.lock.Unlock()
	}

	sg.Unlock(j)

	handleInclude(rc, spec, filename, includes)
}

func handleInclude(
	rc dukkha.ConfigResolvingContext,
	spec *ReadSpec,
	currentFile string,
	include []*IncludeEntry,
) {
	var sg SyncGroup
	sg.Init()

	for i, inc := range include {
		switch {
		case len(inc.Path) != 0:
			toInclude := inc.Path
			if !path.IsAbs(toInclude) {
				// TODO: relative to current file or DUKKHA_WORKDIR ?
				toInclude = path.Join(path.Dir(currentFile), toInclude)
			}

			matches, err := doublestar.Glob(spec.ConfFS, toInclude)
			if err != nil {
				matches = []string{toInclude}
			}

			err = Read(rc, spec, &sg, matches, false)
			if err != nil {
				// TODO: log err?
				_ = fmt.Errorf("loading included config files: %w", err)
				return
			}
		case len(inc.Text) != 0:
			var rd strings.Reader
			rd.Reset(inc.Text)

			go loadConfig(
				rc,
				spec,
				&sg,
				io.NopCloser(&rd),
				fmt.Sprintf("text#%d of %s", i, currentFile),
				sg.NewJob(),
			)
		}
	}

	return
}
