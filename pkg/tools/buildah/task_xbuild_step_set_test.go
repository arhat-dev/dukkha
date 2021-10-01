package buildah

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStepSet_genCtx(t *testing.T) {
	vTrue := true
	pTrue := &vTrue
	vStr := "foo"
	pStr := &vStr

	tests := []struct {
		name     string
		spec     *stepSet
		expected *xbuildContext
	}{
		{
			name: "Commit",
			spec: &stepSet{
				Commit: pTrue,
			},
			expected: &xbuildContext{
				Commit: true,
			},
		},
		{
			name: "Workdir",
			spec: &stepSet{
				Workdir: pStr,
			},
			expected: &xbuildContext{
				WorkDir: vStr,
			},
		},
		{
			name: "User",
			spec: &stepSet{
				User: pStr,
			},
			expected: &xbuildContext{
				User: vStr,
			},
		},
		{
			name: "Shell",
			spec: &stepSet{
				Shell: []string{"sh", "-c", "foo"},
			},
			expected: &xbuildContext{
				Shell: []string{"sh", "-c", "foo"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.EqualValues(t, test.expected, test.spec.genCtx(&xbuildContext{}))
		})
	}
}
