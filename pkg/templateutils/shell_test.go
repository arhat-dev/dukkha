package templateutils

import (
	"bytes"
	"context"
	"encoding/hex"
	"io"
	"strings"
	"testing"

	"arhat.dev/pkg/md5helper"
	"github.com/stretchr/testify/assert"
	"mvdan.cc/sh/v3/syntax"

	di "arhat.dev/dukkha/internal"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
)

func TestEmbeddedShellForTemplateFunc(t *testing.T) {
	tests := []struct {
		name     string
		script   string
		expected string
	}{
		{
			name:     "Simple md5sum",
			script:   `tpl:md5 \"test\"`,
			expected: hex.EncodeToString(md5helper.Sum([]byte("test"))),
		},
		{
			name:     "Piped md5sum",
			script:   `printf "test" | tpl:md5`,
			expected: hex.EncodeToString(md5helper.Sum([]byte("test"))),
		},
		{
			name:     "Subcmd md5sum",
			script:   `tpl:md5 \"$(printf "test")\"`,
			expected: hex.EncodeToString(md5helper.Sum([]byte("test"))),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := dt.NewTestContext(context.TODO())
			ctx.(di.CacheDirSetter).SetCacheDir(t.TempDir())

			stdin, _ := io.Pipe()
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			runner, err := CreateEmbeddedShellRunner("", ctx, stdin, stdout, stderr)
			if !assert.NoError(t, err) {
				return
			}
			assert.Zero(t, stderr.Len())
			stdout.Reset()

			parser := syntax.NewParser(syntax.Variant(syntax.LangBash))
			assert.NoError(t, RunScriptInEmbeddedShell(ctx, runner, parser, test.script))

			assert.Equal(t, test.expected, stdout.String())
			assert.EqualValues(t, 0, stderr.Len())
		})
	}
}

func TestExecCmdAsTemplateFuncCall(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		args      []string
		expected  string
		expectErr bool
	}{
		{
			name:      "Invalid Empty",
			expectErr: true,
		},
		{
			name:     "Valid Simple md5sum",
			args:     []string{"md5", `"test"`},
			expected: hex.EncodeToString(md5helper.Sum([]byte("test"))),
		},
		{
			name:      "Invalid Template Func Not Defined",
			args:      []string{"NOT_DEFINED"},
			expectErr: true,
		},
		{
			name:     "Input for md5sum",
			input:    "test",
			args:     []string{"md5"},
			expected: hex.EncodeToString(md5helper.Sum([]byte("test"))),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			ctx := dt.NewTestContext(context.TODO())
			ctx.(di.CacheDirSetter).SetCacheDir(t.TempDir())

			var input io.Reader
			if len(test.input) != 0 {
				input = strings.NewReader(test.input)
			}

			err := ExecCmdAsTemplateFuncCall(ctx, input, buf, test.args)
			if test.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expected, buf.String())
		})
	}
}
