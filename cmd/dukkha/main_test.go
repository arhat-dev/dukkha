package main

import (
	"os"
	"testing"

	"arhat.dev/dukkha/pkg/sliceutils"
)

func TestMain(t *testing.T) {
	_ = t
	for i, arg := range os.Args {
		if arg == "--" {
			os.Args = sliceutils.NewStrings(
				os.Args[:1], os.Args[i+1:]...,
			)
			break
		}
	}

	main()
}
