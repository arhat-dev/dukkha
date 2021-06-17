package shell_file

import (
	"fmt"

	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/shell"
)

func init() {
	renderer.Register(&Config{}, NewDriver)
}

func NewDriver(config interface{}) (renderer.Interface, error) {
	cfg, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("unexpected non %s renderer config: %T", DefaultName, config)
	}

	if cfg.ExecFunc == nil {
		return nil, fmt.Errorf("required exec func not set")
	}

	return shell.NewDriver(&shell.Config{ExecFunc: cfg.ExecFunc})
}
