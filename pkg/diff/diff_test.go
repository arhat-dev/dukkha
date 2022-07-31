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
		func() *TestCase { return new(TestCase) },
		func() *[]*Expected {
			var keys []*Expected
			return &keys
		},
		func(t *testing.T, in *TestCase, exp *[]*Expected) {
			var actualEntries []*Expected
			for _, ent := range Diff(in.Original, in.Current) {
				actualEntries = append(actualEntries, &Expected{
					Key:       ent.Key,
					Kind:      ent.Kind,
					DivertKey: ent.DivertAt.elemKey,
				})
			}

			if !assert.EqualValues(t, *exp, actualEntries) {
				for i, expected := range *exp {
					assert.EqualValues(t, expected, actualEntries[i])
				}
			}
		},
	)
}
