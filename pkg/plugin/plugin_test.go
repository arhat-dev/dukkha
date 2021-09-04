package plugin

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "embed"
)

// nolint:deadcode,unused,varcheck
//go:embed plugin_example_test.go
var exampleSourceFile string

func TestRef(t *testing.T) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", "\n\nfunc A() string", parser.PackageClauseOnly)
	assert.NoError(t, err)

	assert.Equal(t, "plugin", f.Name.Name)
}
