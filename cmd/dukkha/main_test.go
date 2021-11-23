package main

import (
	"os"
	"testing"

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

	// os.Args = append(os.Args, "run", "workflow", "local", "run", "test")
	main()
}
