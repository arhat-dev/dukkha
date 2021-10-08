package buildah

import (
	"encoding/json"
	"fmt"
	"strings"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
)

type stepSet struct {
	rs.BaseField

	// Workdir
	Workdir *string `yaml:"workdir"`

	// User for command running in run step
	User *string `yaml:"user"`

	// Shell command to interpreter scripts in run step
	Shell []string `yaml:"shell"`

	Env         []*dukkha.EnvEntry `yaml:"env"`
	Annotations []*dukkha.EnvEntry `yaml:"annotations"`
	Labels      []*dukkha.EnvEntry `yaml:"labels"`

	Ports      []string `yaml:"ports"`
	Entrypoint []string `yaml:"entrypoint"`
	Cmd        []string `yaml:"cmd"`
	Volumes    []string `yaml:"volumes"`
	StopSignal *string  `yaml:"stop_signal"`
}

func (s *stepSet) genSpec(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
	record bool,
) ([]dukkha.TaskExecSpec, error) {
	_ = rc

	configCmd := sliceutils.NewStrings(options.ToolCmd(), "config")

	if record {
		configCmd = append(configCmd, "--add-history")
	}

	if s.Workdir != nil {
		configCmd = append(configCmd, "--workingdir", *s.Workdir)
	}

	if s.User != nil {
		configCmd = append(configCmd, "--user", *s.User)
	}

	if len(s.Shell) != 0 {
		configCmd = append(configCmd, "--shell", strings.Join(s.Shell, " "))
	}

	configCmd = append(configCmd, kvArgs("--env", s.Env)...)
	configCmd = append(configCmd, kvArgs("--annotation", s.Annotations)...)
	configCmd = append(configCmd, kvArgs("--label", s.Labels)...)

	for _, p := range s.Ports {
		configCmd = append(configCmd, "--port", p)
	}

	for _, v := range s.Volumes {
		configCmd = append(configCmd, "--volume", v)
	}

	if s.StopSignal != nil {
		configCmd = append(configCmd, "--stop-signal", *s.StopSignal)
	}

	if len(s.Entrypoint) != 0 {
		ent, err := json.Marshal(s.Entrypoint)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal entrypoint value: %w", err)
		}

		configCmd = append(configCmd, "--entrypoint", string(ent))
	}

	if len(s.Cmd) != 0 {
		configCmd = append(configCmd, "--cmd", strings.Join(s.Cmd, " "))
	}

	var steps []dukkha.TaskExecSpec

	steps = append(steps, dukkha.TaskExecSpec{
		IgnoreError: false,
		Command:     append(configCmd, replace_XBUILD_CURRENT_CONTAINER_ID),
		UseShell:    options.UseShell(),
		ShellName:   options.ShellName(),
	})

	return steps, nil
}

func kvArgs(flag string, entries []*dukkha.EnvEntry) []string {
	var ret []string
	for _, a := range entries {
		parts := []string{a.Name}
		if len(a.Value) != 0 {
			parts = append(parts, a.Value)
		}

		ret = append(ret, flag, strings.Join(parts, "="))
	}

	return ret
}
