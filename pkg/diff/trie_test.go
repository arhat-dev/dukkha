package diff

import (
	"testing"

	"arhat.dev/pkg/testhelper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	_ "arhat.dev/rs" // add required references for go:linkname during testing
)

var (
	_ yaml.Marshaler   = (*Node)(nil)
	_ yaml.Unmarshaler = (*Node)(nil)
)

func TestNode_Unmarshal(t *testing.T) {
	type CheckSpec struct {
		Key []string `yaml:"key"`
		// Value is the expected yaml.Node.Value
		Value *string `yaml:"value"`

		Nearest *string `yaml:"nearest"`
	}

	testhelper.TestFixtures(t, "./testdata/trie",
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
					node, exact := n.Get(spec.Key)

					defer func() {
						// kept for debugging purpose
						node, exact := node, exact
						_, _ = node, exact
						if p := recover(); p != nil {
							assert.Failf(t, "",
								"panic for spec %v (size=%d): %v",
								spec.Key, len(spec.Key), p,
							)
						}
					}()

					switch {
					case spec.Nearest != nil:
						assert.False(t, exact)

						assert.NotNil(t, node)
						assert.Equal(t, *spec.Nearest, node.scalarData.Value)
					case spec.Value != nil:
						assert.True(t, exact,
							"expecting exact match of key %v (size=%d), got %v",
							spec.Key, len(spec.Key), node,
						)

						assert.Equal(t, *spec.Value, node.scalarData.Value)
					default:
						assert.Nil(t, node.scalarData)
					}
				}()
			}
		},
	)
}
