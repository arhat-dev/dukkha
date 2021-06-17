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

# build
image.build.dukkha.linux.x86:
	sh scripts/image/build.sh $@

image.build.dukkha.linux.amd64:
	sh scripts/image/build.sh $@

image.build.dukkha.linux.armv5:
	sh scripts/image/build.sh $@

image.build.dukkha.linux.armv6:
	sh scripts/image/build.sh $@

image.build.dukkha.linux.armv7:
	sh scripts/image/build.sh $@

image.build.dukkha.linux.arm64:
	sh scripts/image/build.sh $@

image.build.dukkha.linux.ppc64le:
	sh scripts/image/build.sh $@

image.build.dukkha.linux.mips64le:
	sh scripts/image/build.sh $@

image.build.dukkha.linux.s390x:
	sh scripts/image/build.sh $@

image.build.dukkha.linux.all: \
	image.build.dukkha.linux.amd64 \
	image.build.dukkha.linux.arm64 \
	image.build.dukkha.linux.armv7 \
	image.build.dukkha.linux.armv6 \
	image.build.dukkha.linux.armv5 \
	image.build.dukkha.linux.x86 \
	image.build.dukkha.linux.s390x \
	image.build.dukkha.linux.ppc64le \
	image.build.dukkha.linux.mips64le

image.build.dukkha.windows.amd64:
	sh scripts/image/build.sh $@

image.build.dukkha.windows.armv7:
	sh scripts/image/build.sh $@

image.build.dukkha.windows.all: \
	image.build.dukkha.windows.amd64 \
	image.build.dukkha.windows.armv7

# push
image.push.dukkha.linux.x86:
	sh scripts/image/push.sh $@

image.push.dukkha.linux.amd64:
	sh scripts/image/push.sh $@

image.push.dukkha.linux.armv5:
	sh scripts/image/push.sh $@

image.push.dukkha.linux.armv6:
	sh scripts/image/push.sh $@

image.push.dukkha.linux.armv7:
	sh scripts/image/push.sh $@

image.push.dukkha.linux.arm64:
	sh scripts/image/push.sh $@

image.push.dukkha.linux.ppc64le:
	sh scripts/image/push.sh $@

image.push.dukkha.linux.mips64le:
	sh scripts/image/push.sh $@

image.push.dukkha.linux.s390x:
	sh scripts/image/push.sh $@

image.push.dukkha.linux.all: \
	image.push.dukkha.linux.amd64 \
	image.push.dukkha.linux.arm64 \
	image.push.dukkha.linux.armv7 \
	image.push.dukkha.linux.armv6 \
	image.push.dukkha.linux.armv5 \
	image.push.dukkha.linux.x86 \
	image.push.dukkha.linux.s390x \
	image.push.dukkha.linux.ppc64le \
	image.push.dukkha.linux.mips64le

image.push.dukkha.windows.amd64:
	sh scripts/image/push.sh $@

image.push.dukkha.windows.armv7:
	sh scripts/image/push.sh $@

image.push.dukkha.windows.all: \
	image.push.dukkha.windows.amd64 \
	image.push.dukkha.windows.armv7
