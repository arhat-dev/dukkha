package main

import (
	"os"
	"path/filepath"
	"testing"

	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/sliceutils"
	"github.com/stretchr/testify/assert"
)

func clearOSArgs() {
	for i, arg := range os.Args {
		if arg == "--" {
			os.Args = sliceutils.NewStrings(
				os.Args[:1], os.Args[i+1:]...,
			)
			break
		}
	}
}

func TestMain(t *testing.T) {
	_ = t

	clearOSArgs()
	os.Args = append(os.Args, "workflow", "local", "run", "test")
	main()
}

func TestRender(t *testing.T) {
	_ = t

	clearOSArgs()

	testCases := matrix.CartesianProduct(map[string][]string{
		"format": {"json", "yaml"},
	})

	const (
		outputDir   = "build/cmd-render-test"
		srcDir      = "cmd/dukkha/testdata/render/source"
		expectedDir = "cmd/dukkha/testdata/render/expected"
	)
	baseArgs := sliceutils.NewStrings(
		os.Args, "render", "-o", outputDir, srcDir,
	)

	if !assert.NoError(t, os.RemoveAll(outputDir), "output dir not cleandup") {
		return
	}

	for _, test := range testCases {
		os.Args = sliceutils.NewStrings(baseArgs, "-f", test["format"])
		main()
	}

	entries, err := os.ReadDir(expectedDir)
	if !assert.NoError(t, err) {
		return
	}

	for _, ent := range entries {
		data, err := os.ReadFile(filepath.Join(outputDir, ent.Name()))
		if !assert.NoError(t, err, "test result not generated") {
			return
		}

		expected, err := os.ReadFile(filepath.Join(expectedDir, ent.Name()))
		if !assert.NoError(t, err, "expected data not loaded") {
			return
		}

		assert.EqualValues(t, string(expected), string(data))
	}
}
