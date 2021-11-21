package archive

import (
	"context"
	"strings"
	"testing"

	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"
)

func TestLcpp(t *testing.T) {
	for _, test := range []struct {
		list     []string
		expected string
	}{
		{
			list:     []string{},
			expected: "",
		},
		{
			list:     []string{"a"},
			expected: "",
		},
		{
			list:     []string{"a/"},
			expected: "a/",
		},
		{
			list:     []string{"a/b/c", "a/b/d", "a/b/a"},
			expected: "a/b/",
		},
		{
			list:     []string{"a/x/c", "a/xx/d", "a/xx/a"},
			expected: "a/",
		},
		{
			list:     []string{"a/🌶️🌶️/c", "a/🌶️x/d", "a/🌶️🌶️/a"},
			expected: "a/",
		},
		{
			list:     []string{"a/🌶️🌶️/c", "a/🌶️🌶️/d", "a/🌶️🌶️/a"},
			expected: "a/🌶️🌶️/",
		},
	} {
		assert.EqualValues(t, test.expected, lcpp(test.list))
	}
}

func TestCollectFiles(t *testing.T) {
	type TestCase struct {
		rs.BaseField

		Task      *TaskCreate `yaml:"task"`
		ExpectErr bool        `yaml:"expect_err"`
	}

	type ExpectedEntry struct {
		From string `yaml:"from"`
		Link string `yaml:"link"`
	}

	testhelper.TestFixtures(t, "./fixtures/collect-files",
		func() interface{} { return rs.Init(&TestCase{}, nil) },
		func() interface{} {
			m := make(map[string]*ExpectedEntry)
			return &m
		},
		func(t *testing.T, in, exp interface{}) {
			spec := in.(*TestCase)
			expected := *exp.(*map[string]*ExpectedEntry)

			ctx := dukkha_test.NewTestContext(context.TODO())
			assert.NoError(t, spec.ResolveFields(ctx, -1))

			actualFiles, err := collectFiles(spec.Task.Files)
			if spec.ExpectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			assert.EqualValues(t, len(expected), len(actualFiles))

			var files []string
			m := make(map[string]*entry)
			for _, v := range actualFiles {
				m[v.to] = v
				files = append(files, v.to)
			}

			t.Log(strings.Join(files, ", "))

			for k, exp := range expected {
				t.Run(k, func(t *testing.T) {
					actual, ok := m[k]
					if !assert.True(t, ok, "%q not found", k) {
						return
					}

					assert.EqualValues(t, exp.From, actual.from, "bad source")
					assert.EqualValues(t, exp.Link, actual.link, "bad link")
				})
			}
		},
	)
}
