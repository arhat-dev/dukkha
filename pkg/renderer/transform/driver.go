package transform

import (
	"fmt"

	"arhat.dev/pkg/stringhelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer"
)

const (
	DefaultName = "T"
)

func init() { dukkha.RegisterRenderer(DefaultName, NewDefault) }

func NewDefault(name string) dukkha.Renderer { return &Driver{name: name} }

type Driver struct {
	rs.BaseField `yaml:"-"`

	renderer.BaseRenderer `yaml:",inline"`

	name string
}

func (d *Driver) RenderYaml(
	rc dukkha.RenderingContext, rawData interface{}, _ []dukkha.RendererAttribute,
) ([]byte, error) {
	rawData, err := rs.NormalizeRawData(rawData)
	if err != nil {
		return nil, err
	}

	rawBytes, err := yamlhelper.ToYamlBytes(rawData)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: unsupported input type %T: %w",
			d.name, rawData, err,
		)
	}

	spec := rs.Init(&Spec{}, nil).(*Spec)
	err = yaml.Unmarshal(rawBytes, spec)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: unmarshal transform spec: %w",
			d.name, err,
		)
	}

	// only resolve value and ops list, we need to resolve
	// each operation step by step with `value` injected
	err = spec.ResolveFields(rc, 2)
	if err != nil {
		return nil, fmt.Errorf(
			"renderer.%s: resolving transform spec: %w",
			d.name, err,
		)
	}

	data := spec.Value
	for i, op := range spec.Ops {
		data, err = op.Do(rc, data)
		if err != nil {
			return nil, fmt.Errorf(
				"renderer.%s: step #%d: %w",
				d.name, i, err,
			)
		}
	}

	return stringhelper.ToBytes[byte, byte](data), nil
}

// Spec for yaml data transformation
type Spec struct {
	rs.BaseField `yaml:"-"`

	// Value always be string, so you can decide which format is converts to in operations
	Value string `yaml:"value"`

	// Ops the transform operations to run
	Ops []*Operation `yaml:"ops"`
}

type Operation struct {
	rs.BaseField `yaml:"-"`

	// AWK runs an awk script to process VALUE as input
	AWK *awkSpec `yaml:"awk,omitempty"`

	// TLang runs a tlang script to process VALUE
	TLang *tlangSpec `yaml:"tlang,omitempty"`

	// Tmpl executes a golang template to process VALUE
	Tmpl *tmplSpec `yaml:"tmpl,omitempty"`

	// Shell runs a bash script in the embedded bash shell to process VALUE
	Shell *shellSpec `yaml:"shell,omitempty"`

	// Checksum verifies checksum of the VALUE and leaves VALUE unchanged
	Checksum *Checksum `yaml:"checksum,omitempty"`
}

type extendedUserFacingRenderContext interface {
	dukkha.RenderingContext
	di.VALUEGetter
}

func (op *Operation) Do(_rc dukkha.RenderingContext, value string) (_ string, err error) {
	rc2 := _rc.(interface {
		dukkha.RenderingContext
		di.VALUEGetter
		di.VALUESetter
	})

	rc2.SetVALUE(value)
	rc2.AddEnv(true, &dukkha.NameValueEntry{
		Name:  "VALUE",
		Value: rc2.VALUE().(string),
	})

	// do not expose SetVALUE to template operation
	rc := rc2.(extendedUserFacingRenderContext)

	err = op.ResolveFields(rc, -1)
	if err != nil {
		return
	}

	switch {
	case op.AWK != nil:
		return op.AWK.Run(rc, value)
	case op.TLang != nil:
		return op.TLang.Run(rc)
	case op.Tmpl != nil:
		return op.Tmpl.Run(rc)
	case op.Shell != nil:
		return op.Shell.Run(rc)
	case op.Checksum != nil:
		return value, op.Checksum.Verify(rc.FS())
	default:
		return value, nil
	}
}
