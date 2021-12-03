package conf

import (
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/rs"
)

type GlobalConfig struct {
	rs.BaseField `yaml:"-"`

	// CacheDir set DUKKHA_CACHE_DIR to store script files, renderer cache
	// and intermediate task execution data
	CacheDir string `yaml:"cache_dir"`

	// DefaultGitBranch set GIT_DEFAULT_BRANCH, useful when dukkha can not
	// detect branch name of origin/HEAD (e.g. github ci environment)
	//
	// If your have multiple definitions of this option in different config
	// file, only the first occurrence of the option is used.
	DefaultGitBranch string `yaml:"default_git_branch"`

	// Env add global environment variables for all working parts in dukkha
	Env dukkha.Env `yaml:"env"`

	// Values is the global store of runtime values
	//
	// accessible from renderer template `{{ values.YOUR_VAL_KEY }}`
	// and renderer env/shell `${values.YOUR_VAL_KEY}`
	Values rs.AnyObjectMap `yaml:"values"`
}

func (g *GlobalConfig) Merge(a *GlobalConfig) error {
	err := g.BaseField.Inherit(&a.BaseField)
	if err != nil {
		return fmt.Errorf("failed to inherit other global config: %w", err)
	}

	g.Env = append(g.Env, a.Env...)
	if len(a.CacheDir) != 0 {
		g.CacheDir = a.CacheDir
	}

	if len(a.DefaultGitBranch) != 0 {
		g.DefaultGitBranch = a.DefaultGitBranch
	}

	err = g.Values.Inherit(&a.Values.BaseField)
	if err != nil {
		return fmt.Errorf("failed to merge global values: %w", err)
	}

	if len(a.Values.Data) != 0 {
		if g.Values.Data == nil {
			g.Values.Data = a.Values.Data
		} else {
			for k, v := range a.Values.Data {
				g.Values.Data[k] = v
			}
		}
	}

	return nil
}

// ResolveAllButValues resolves global env first, and make them available to
// all other fields, then resolves all other fields except values (will be handled later)
func (g *GlobalConfig) ResolveAllButValues(rc dukkha.ConfigResolvingContext) error {
	err := dukkha.ResolveEnv(rc, g, "Env", "env")
	if err != nil {
		return fmt.Errorf("failed to resolve global env: %w", err)
	}

	err = g.ResolveFields(rc, -1, "cache_dir")
	if err != nil {
		return fmt.Errorf("failed to resolve cache dir: %w", err)
	}

	err = g.ResolveFields(rc, -1, "default_git_branch")
	if err != nil {
		return fmt.Errorf("failed to resolve default git branch: %w", err)
	}

	return nil
}
