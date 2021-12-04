package input

import (
	"fmt"
	"os"

	"arhat.dev/pkg/iohelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
)

const (
	DefaultName = "input"
)

func init() { dukkha.RegisterRenderer(DefaultName, NewDefault) }

func NewDefault(name string) dukkha.Renderer { return &Driver{name: name} }

var _ dukkha.Renderer = (*Driver)(nil)

type Driver struct {
	rs.BaseField `yaml:"-"`

	renderer.BaseRenderer `yaml:",inline"`

	name string

	Config configSpec `yaml:",inline"`
}

func (d *Driver) RenderYaml(
	_ dukkha.RenderingContext, rawData interface{}, attributes []dukkha.RendererAttribute,
) ([]byte, error) {
	rawData, err := rs.NormalizeRawData(rawData)
	if err != nil {
		return nil, err
	}

	promptBytes, err := yamlhelper.ToYamlBytes(rawData)
	if err != nil {
		return nil, err
	}

	var useSpec bool
	for _, attr := range d.Attributes(attributes) {
		switch attr {
		case renderer.AttrUseSpec:
			useSpec = true
		default:
		}
	}

	var (
		prompt = d.Config.Prompt
		hide   = d.Config.HideInput
	)
	if useSpec {
		spec := rs.Init(&inputSpec{}, nil).(*inputSpec)
		err = yaml.Unmarshal(promptBytes, spec)
		if err != nil {
			return nil, fmt.Errorf("renderer.%s: invalid input spec %w", d.name, err)
		}

		if spec.Config.HideInput != nil {
			hide = spec.Config.HideInput
		}

		if len(spec.Config.Prompt) != 0 {
			prompt = spec.Config.Prompt
		}
	} else if len(promptBytes) != 0 {
		prompt = string(promptBytes)
	}

	fmt.Print(prompt)

	if hide != nil && *hide {
		ret, err2 := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		return ret, err2
	}

	ret, err := iohelper.ReadInputLine(os.Stdin)
	fmt.Println()
	return ret, err
}
