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

func init() { dukkha.RegisterRenderer(DefaultName, NewDefault) }

func NewDefault(name string) dukkha.Renderer { return &Driver{name: name} }

var _ dukkha.Renderer = (*Driver)(nil)

type Driver struct {
	rs.BaseField `yaml:"-"`

	Hide bool `yaml:"hide"`

	name string
}

func (d *Driver) Init(ctx dukkha.ConfigResolvingContext) error {
	return nil
}

func (d *Driver) RenderYaml(
	_ dukkha.RenderingContext, rawData interface{},
) ([]byte, error) {
	rawData, err := rs.NormalizeRawData(rawData)
	if err != nil {
		return nil, err
	}

	promptBytes, err := yamlhelper.ToYamlBytes(rawData)
	if err != nil {
		return nil, err
	}

	fmt.Print(string(promptBytes))
	if d.Hide {
		ret, err2 := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		return ret, err2
	}

	ret, err := readline(os.Stdin)
	fmt.Println()
	return ret, err
}
