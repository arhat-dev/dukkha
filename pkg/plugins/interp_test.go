package plugins

import (
	"context"
	"testing"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/plugins/internal"
)

const testSrcFib = `package main

func main() {
	fib(35)
}

func fib(i int) int {
	switch i {
	case 0:
		return 0
	case 1:
		return 1
	}

	return fib(i-1) + fib(i-2)
}
`

func TestInterp(t *testing.T) {
	buildOpts := scriggo.BuildOptions{
		AllowGoStmt: false,
		Packages:    internal.NativePackages,
		Globals:     native.Declarations{},
	}

	runOpts := scriggo.RunOptions{
		Context: context.TODO(),
		Print: func(i any) {

		},
	}

	prog, err := scriggo.Build(scriggo.Files{
		"main.go": []byte(`package main

import "fmt"

func main() {
	fmt.Println("hello")
}
`),
	}, &buildOpts)
	assert.NoError(t, err)

	assert.NoError(t, prog.Run(&runOpts))
}

func BenchmarkBuild(b *testing.B) {
	buildOpts := scriggo.BuildOptions{
		AllowGoStmt: false,
		Packages:    internal.NativePackages,
		Globals:     native.Declarations{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = scriggo.Build(scriggo.Files{
			"main.go": []byte(testSrcFib),
		}, &buildOpts)
	}
}

func BenchmarkRun(b *testing.B) {
	buildOpts := scriggo.BuildOptions{
		AllowGoStmt: false,
		Packages:    internal.NativePackages,
		Globals:     native.Declarations{},
	}

	prog, err := scriggo.Build(scriggo.Files{
		"main.go": []byte(testSrcFib),
	}, &buildOpts)
	assert.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = prog.Run(nil)
	}
}
