//go:build !darwin && !linux && !freebsd && !openbsd && !netbsd && !dragonfly
// +build !darwin,!linux,!freebsd,!openbsd,!netbsd,!dragonfly

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

package iohelper

import (
	"arhat.dev/pkg/errhelper"
)

const errNOSYS errhelper.ErrString = "function not implemented"

func CheckBytesToRead(fd uintptr) (int, error) {
	return 0, errNOSYS
}
