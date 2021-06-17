package template

import (
	"fmt"

	"arhat.dev/dukkha/pkg/renderer"
)

func init() {
	renderer.Register(&Config{}, NewDriver)
}

func NewDriver(config interface{}) (renderer.Interface, error) {
	cfg, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("unexpected non template renderer config: %T", config)
	}

	_ = cfg

	return &Driver{}, nil
}
