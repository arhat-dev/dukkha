package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/sliceutils"
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
	os.Args = append(os.Args, "run", "workflow", "local", "run", "test")
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
		os.Args, "render", "-q", ".", "-o", outputDir, srcDir,
	)

	if !assert.NoError(t, os.RemoveAll(outputDir), "output dir not cleandup") {
		return
	}

	wd, err := os.Getwd()
	if !assert.NoError(t, err, "unable to get working dir") {
		return
	}

	if !assert.True(t, filepath.IsAbs(wd), "working dir is not absolute path") {
		return
	}

	for _, test := range testCases {
		os.Args = sliceutils.NewStrings(baseArgs, "-f", test["format"])
		main()
		if !assert.NoError(t, os.Chdir(wd)) {
			assert.FailNow(t, "unable to go back to original working dir")
		}
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
