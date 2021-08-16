package input

import (
	"fmt"
	"os"

	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"golang.org/x/term"

	"arhat.dev/dukkha/pkg/dukkha"
)

// nolint:revive
const (
	DefaultName = "input"
)

func init() {
	dukkha.RegisterRenderer(DefaultName, NewDefault)
}

func NewDefault() dukkha.Renderer {
	return &driver{}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	rs.BaseField
}

func (d *driver) Init(ctx dukkha.ConfigResolvingContext) error {
	ctx.AddRenderer(DefaultName, d)
	return nil
}

func (d *driver) RenderYaml(_ dukkha.RenderingContext, rawData interface{}) ([]byte, error) {
	promptBytes, err := yamlhelper.ToYamlBytes(rawData)
	if err != nil {
		return nil, err
	}

	fmt.Print(string(promptBytes))
	ret, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	return ret, err
}
