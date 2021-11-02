# Copyright 2020 The arhat.dev Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

include scripts/lint.mk

GOMOD := GOPROXY=direct GOSUMDB=off go mod
.PHONY: vendor
vendor:
	${GOMOD} tidy
	${GOMOD} vendor
	patch -u -p1 --verbose -i scripts/patches/vendor.patch

GOOS ?= $(shell go env GOHOSTOS)
GOARCH ?= $(shell go env GOHOSTARCH)
.PHONY: dukkha
dukkha:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} \
		go build -o build/dukkha -mod vendor ./cmd/dukkha

.PHONY: docs
docs:
	go test -v -mod=readonly -tags="docs" ./docs

# testing
include scripts/test/unit.mk

# packaging
include scripts/package/dukkha.mk

# optional private scripts
-include private/scripts.mk
