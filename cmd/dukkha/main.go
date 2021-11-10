/*
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"arhat.dev/dukkha/pkg/cmd"
	"arhat.dev/dukkha/pkg/version"

	_ "arhat.dev/dukkha/cmd/dukkha/addon"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(version.NewVersionCmd())

	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())

		exitErr := new(exec.ExitError)
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}

		os.Exit(1)
	}
}
