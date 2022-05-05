package diff

import (
	"testing"

	"arhat.dev/pkg/testhelper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var (
	_ yaml.Marshaler   = (*Node)(nil)
	_ yaml.Unmarshaler = (*Node)(nil)
)

func TestNode_Unmarshal(t *testing.T) {
	t.Parallel()

	type CheckSpec struct {
		Key []string `yaml:"key"`
		// Value is the expected yaml.Node.Value
		Value *string `yaml:"value"`

		Nearest *string `yaml:"nearest"`

		TailKey []string `yaml:"tail_key"`
	}

	testhelper.TestFixtures(t, "./fixtures/trie",
		func() interface{} { return new(Node) },
		func() interface{} {
			var specs []CheckSpec
			return &specs
		},
		func(t *testing.T, in, cs interface{}) {
			n := in.(*Node)
			checkSpecs := *cs.(*[]CheckSpec)

			for _, spec := range checkSpecs {
				func() {
					node, tailKey := n.Get(spec.Key)

					defer func() {
						// kept for debugging purpose
						_, _ = node, tailKey
						if p := recover(); p != nil {
							assert.Failf(t, "",
								"panic for spec %v (size=%d): %v",
								spec.Key, len(spec.Key), p,
							)
						}
					}()

					assert.EqualValues(t, spec.TailKey, tailKey)

					switch {
					case spec.Nearest != nil:
						assert.NotNil(t, node)
						assert.Equal(t, *spec.Nearest, node.scalarData.Value)
					case spec.Value != nil:
						assert.Equal(t, *spec.Value, node.scalarData.Value)
					default:
						assert.Nil(t, node.scalarData)
					}
				}()
			}
		},
	)
}
