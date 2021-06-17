package conf

import "context"

type ShellConfig struct {
	Name    Field `dukkha:"name,string"`
	Path    Field `dukkha:"path,string"`
	Env     Field `dukkha:"field,[]string"`
	Command Field `dukkha:"command,[]string"`
	Args    Field `dukkha:"args,[]string"`
}

type ShellConfigList []ShellConfig

func (c *ShellConfigList) resolve(ctx context.Context, data interface{}) error {
	return nil
}
