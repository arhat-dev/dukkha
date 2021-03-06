package buildah

import (
	"encoding/json"
	"fmt"
	"strings"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

type stepSet struct {
	rs.BaseField `yaml:"-"`

	// Workdir
	Workdir *string `yaml:"workdir"`

	// User for command running in run step
	User *string `yaml:"user"`

	// Shell command to interpreter scripts in run step
	Shell []string `yaml:"shell"`

	Env         []*dukkha.NameValueEntry `yaml:"env"`
	Annotations []*dukkha.NameValueEntry `yaml:"annotations"`
	Labels      []*dukkha.NameValueEntry `yaml:"labels"`

	Ports      []string `yaml:"ports"`
	Entrypoint []string `yaml:"entrypoint"`
	Cmd        []string `yaml:"cmd"`
	Volumes    []string `yaml:"volumes"`
	StopSignal *string  `yaml:"stop_signal"`
}

func (s *stepSet) genSpec(
	rc dukkha.TaskExecContext,
	record bool,
) ([]dukkha.TaskExecSpec, error) {
	_ = rc

	var configArgs []string

	if s.Workdir != nil {
		configArgs = append(configArgs, "--workingdir", *s.Workdir)
	}

	if s.User != nil {
		configArgs = append(configArgs, "--user", *s.User)
	}

	if len(s.Shell) != 0 {
		configArgs = append(configArgs, "--shell", strings.Join(s.Shell, " "))
	}

	configArgs = append(configArgs, kvArgs("--env", s.Env)...)
	configArgs = append(configArgs, kvArgs("--annotation", s.Annotations)...)
	configArgs = append(configArgs, kvArgs("--label", s.Labels)...)

	for _, p := range s.Ports {
		configArgs = append(configArgs, "--port", p)
	}

	for _, v := range s.Volumes {
		configArgs = append(configArgs, "--volume", v)
	}

	if s.StopSignal != nil {
		configArgs = append(configArgs, "--stop-signal", *s.StopSignal)
	}

	if len(s.Entrypoint) != 0 {
		ent, err := json.Marshal(s.Entrypoint)
		if err != nil {
			return nil, fmt.Errorf("xbuild: marshal entrypoint: %w", err)
		}

		configArgs = append(configArgs, "--entrypoint", string(ent))
	}

	if len(s.Cmd) != 0 {
		configArgs = append(configArgs, "--cmd", strings.Join(s.Cmd, " "))
	}

	if len(configArgs) == 0 {
		// no config updated
		return nil, nil
	}

	configCmd := []string{constant.DUKKHA_TOOL_CMD, "config"}
	if record {
		configCmd = append(configCmd, "--add-history")
	}

	configCmd = append(configCmd, configArgs...)

	var steps []dukkha.TaskExecSpec

	steps = append(steps, dukkha.TaskExecSpec{
		IgnoreError: false,
		Command:     append(configCmd, replace_XBUILD_CURRENT_CONTAINER_ID),
	})

	return steps, nil
}

func kvArgs(flag string, entries []*dukkha.NameValueEntry) []string {
	var ret []string
	for _, a := range entries {
		parts := []string{a.Name}

		// if annotation/label/env has a trailing dash (`-`),
		// then it is removed from the config
		//
		// we should not set value in that case
		if !strings.HasSuffix(a.Name, "-") {
			parts = append(parts, a.Value)
		}

		ret = append(ret, flag, strings.Join(parts, "="))
	}

	return ret
}
