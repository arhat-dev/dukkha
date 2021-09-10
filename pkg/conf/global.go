package conf

import (
	"fmt"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
)

type GlobalConfig struct {
	rs.BaseField

	// CacheDir to store script file and temporary task execution data
	CacheDir string `yaml:"cache_dir"`

	DefaultGitBranch string `yaml:"default_git_branch"`

	// Env
	Env dukkha.Env `yaml:"env"`

	Values dukkha.ArbitraryValues `yaml:"values"`
}

func (g *GlobalConfig) Merge(a *GlobalConfig) error {
	if a == nil {
		return nil
	}

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

	err = g.Values.ShallowMerge(&a.Values)
	if err != nil {
		return fmt.Errorf("failed to merge global values: %w", err)
	}

	return nil
}

func (g *GlobalConfig) ResolveAllButValues(rc dukkha.ConfigResolvingContext) error {
	err := dukkha.ResolveEnv(g, rc, "Env")
	if err != nil {
		return fmt.Errorf("failed to resolve global env: %w", err)
	}

	err = g.ResolveFields(rc, -1, "CacheDir")
	if err != nil {
		return fmt.Errorf("failed to resolve cache dir: %w", err)
	}

	err = g.ResolveFields(rc, -1, "DefaultGitBranch")
	if err != nil {
		return fmt.Errorf("failed to resolve default git branch: %w", err)
	}

	return nil
}
