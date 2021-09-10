package conf

import (
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/plugin"
	"arhat.dev/pkg/rshelper"
	"arhat.dev/rs"
)

func NewPluginConfig() *PluginConfig {
	return rshelper.InitAll(
		&PluginConfig{},
		dukkha.GlobalInterfaceTypeHandler,
	).(*PluginConfig)
}

type PluginConfig struct {
	rs.BaseField

	Renderers []*plugin.RendererSpec `yaml:"renderers"`
	Tools     []*plugin.ToolSpec     `yaml:"tools"`
	Tasks     []*plugin.TaskSpec     `yaml:"tasks"`
}

func (p *PluginConfig) Merge(a *PluginConfig) error {
	if a == nil {
		return nil
	}

	err := p.BaseField.Inherit(&a.BaseField)
	if err != nil {
		return fmt.Errorf("failed to inherit other plugins config: %w", err)
	}

	p.Renderers = append(p.Renderers, a.Renderers...)
	p.Tools = append(p.Tools, a.Tools...)
	p.Tasks = append(p.Tasks, a.Tasks...)

	return nil
}

func (p *PluginConfig) ResolveAndRegisterPlugins(
	appCtx dukkha.ConfigResolvingContext,
) error {
	err := p.ResolveFields(appCtx, -1)
	if err != nil {
		return fmt.Errorf("failed to resolve renderers: %w", err)
	}

	for i, r := range p.Renderers {
		err = r.Register(p.Renderers[i], appCtx.CacheDir())
		if err != nil {
			return fmt.Errorf("failed to register renderer plugin %q: %w", r.Name(), err)
		}
	}

	for i, t := range p.Tools {
		err = t.Register(p.Tools[i], appCtx.CacheDir())
		if err != nil {
			return fmt.Errorf("failed to register tool plugin %q: %w", t.Name(), err)
		}
	}

	for i, t := range p.Tasks {
		err = t.Register(p.Tasks[i], appCtx.CacheDir())
		if err != nil {
			return fmt.Errorf("failed to register task plugin %q: %w", t.Name(), err)
		}
	}

	return nil
}
