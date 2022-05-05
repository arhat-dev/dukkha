package diff

import (
	"testing"

	"arhat.dev/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		Original *Node `yaml:"original"`
		Current  *Node `yaml:"current"`
	}

	type Expected struct {
		Key       []string `yaml:"key"`
		Kind      Kind     `yaml:"kind"`
		DivertKey string   `yaml:"divert_key"`
	}

	testhelper.TestFixtures(t, "./fixtures/diff",
		func() interface{} { return new(TestCase) },
		func() interface{} {
			var keys []*Expected
			return &keys
		},
		func(t *testing.T, spec, exp interface{}) {
			s := spec.(*TestCase)
			expectedEntries := *exp.(*[]*Expected)

			var actualEntries []*Expected
			for _, ent := range Diff(s.Original, s.Current) {
				actualEntries = append(actualEntries, &Expected{
					Key:       ent.Key,
					Kind:      ent.Kind,
					DivertKey: ent.DivertAt.elemKey,
				})
			}

			if !assert.EqualValues(t, expectedEntries, actualEntries) {
				for i, expected := range expectedEntries {
					assert.EqualValues(t, expected, actualEntries[i])
				}
			}
		},
	)
}
