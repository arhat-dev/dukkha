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

# native
dukkha:
	sh scripts/build/build.sh $@

# linux
dukkha.linux.x86:
	sh scripts/build/build.sh $@

dukkha.linux.amd64:
	sh scripts/build/build.sh $@

dukkha.linux.armv5:
	sh scripts/build/build.sh $@

dukkha.linux.armv6:
	sh scripts/build/build.sh $@

dukkha.linux.armv7:
	sh scripts/build/build.sh $@

dukkha.linux.arm64:
	sh scripts/build/build.sh $@

dukkha.linux.mips:
	sh scripts/build/build.sh $@

dukkha.linux.mipshf:
	sh scripts/build/build.sh $@

dukkha.linux.mipsle:
	sh scripts/build/build.sh $@

dukkha.linux.mipslehf:
	sh scripts/build/build.sh $@

dukkha.linux.mips64:
	sh scripts/build/build.sh $@

dukkha.linux.mips64hf:
	sh scripts/build/build.sh $@

dukkha.linux.mips64le:
	sh scripts/build/build.sh $@

dukkha.linux.mips64lehf:
	sh scripts/build/build.sh $@

dukkha.linux.ppc64:
	sh scripts/build/build.sh $@

dukkha.linux.ppc64le:
	sh scripts/build/build.sh $@

dukkha.linux.s390x:
	sh scripts/build/build.sh $@

dukkha.linux.riscv64:
	sh scripts/build/build.sh $@

dukkha.linux.all: \
	dukkha.linux.x86 \
	dukkha.linux.amd64 \
	dukkha.linux.armv5 \
	dukkha.linux.armv6 \
	dukkha.linux.armv7 \
	dukkha.linux.arm64 \
	dukkha.linux.mips \
	dukkha.linux.mipshf \
	dukkha.linux.mipsle \
	dukkha.linux.mipslehf \
	dukkha.linux.mips64 \
	dukkha.linux.mips64hf \
	dukkha.linux.mips64le \
	dukkha.linux.mips64lehf \
	dukkha.linux.ppc64 \
	dukkha.linux.ppc64le \
	dukkha.linux.s390x \
	dukkha.linux.riscv64

dukkha.darwin.amd64:
	sh scripts/build/build.sh $@

dukkha.darwin.arm64:
	sh scripts/build/build.sh $@

dukkha.darwin.all: \
	dukkha.darwin.amd64

dukkha.windows.x86:
	sh scripts/build/build.sh $@

dukkha.windows.amd64:
	sh scripts/build/build.sh $@

dukkha.windows.armv5:
	sh scripts/build/build.sh $@

dukkha.windows.armv6:
	sh scripts/build/build.sh $@

dukkha.windows.armv7:
	sh scripts/build/build.sh $@

# # currently no support for windows/arm64
# dukkha.windows.arm64:
# 	sh scripts/build/build.sh $@

dukkha.windows.all: \
	dukkha.windows.amd64 \
	dukkha.windows.armv7 \
	dukkha.windows.x86 \
	dukkha.windows.armv5 \
	dukkha.windows.armv6

# # android build requires android sdk
# dukkha.android.amd64:
# 	sh scripts/build/build.sh $@

# dukkha.android.x86:
# 	sh scripts/build/build.sh $@

# dukkha.android.armv5:
# 	sh scripts/build/build.sh $@

# dukkha.android.armv6:
# 	sh scripts/build/build.sh $@

# dukkha.android.armv7:
# 	sh scripts/build/build.sh $@

# dukkha.android.arm64:
# 	sh scripts/build/build.sh $@

# dukkha.android.all: \
# 	dukkha.android.amd64 \
# 	dukkha.android.arm64 \
# 	dukkha.android.x86 \
# 	dukkha.android.armv7 \
# 	dukkha.android.armv5 \
# 	dukkha.android.armv6

dukkha.freebsd.amd64:
	sh scripts/build/build.sh $@

dukkha.freebsd.x86:
	sh scripts/build/build.sh $@

dukkha.freebsd.armv5:
	sh scripts/build/build.sh $@

dukkha.freebsd.armv6:
	sh scripts/build/build.sh $@

dukkha.freebsd.armv7:
	sh scripts/build/build.sh $@

dukkha.freebsd.arm64:
	sh scripts/build/build.sh $@

dukkha.freebsd.all: \
	dukkha.freebsd.amd64 \
	dukkha.freebsd.arm64 \
	dukkha.freebsd.armv7 \
	dukkha.freebsd.x86 \
	dukkha.freebsd.armv5 \
	dukkha.freebsd.armv6

dukkha.netbsd.amd64:
	sh scripts/build/build.sh $@

dukkha.netbsd.x86:
	sh scripts/build/build.sh $@

dukkha.netbsd.armv5:
	sh scripts/build/build.sh $@

dukkha.netbsd.armv6:
	sh scripts/build/build.sh $@

dukkha.netbsd.armv7:
	sh scripts/build/build.sh $@

dukkha.netbsd.arm64:
	sh scripts/build/build.sh $@

dukkha.netbsd.all: \
	dukkha.netbsd.amd64 \
	dukkha.netbsd.arm64 \
	dukkha.netbsd.armv7 \
	dukkha.netbsd.x86 \
	dukkha.netbsd.armv5 \
	dukkha.netbsd.armv6

dukkha.openbsd.amd64:
	sh scripts/build/build.sh $@

dukkha.openbsd.x86:
	sh scripts/build/build.sh $@

dukkha.openbsd.armv5:
	sh scripts/build/build.sh $@

dukkha.openbsd.armv6:
	sh scripts/build/build.sh $@

dukkha.openbsd.armv7:
	sh scripts/build/build.sh $@

dukkha.openbsd.arm64:
	sh scripts/build/build.sh $@

dukkha.openbsd.all: \
	dukkha.openbsd.amd64 \
	dukkha.openbsd.arm64 \
	dukkha.openbsd.armv7 \
	dukkha.openbsd.x86 \
	dukkha.openbsd.armv5 \
	dukkha.openbsd.armv6

dukkha.plan9.amd64:
	sh scripts/build/build.sh $@

dukkha.plan9.x86:
	sh scripts/build/build.sh $@

dukkha.plan9.armv5:
	sh scripts/build/build.sh $@

dukkha.plan9.armv6:
	sh scripts/build/build.sh $@

dukkha.plan9.armv7:
	sh scripts/build/build.sh $@

dukkha.plan9.all: \
	dukkha.plan9.amd64 \
	dukkha.plan9.armv7 \
	dukkha.plan9.x86 \
	dukkha.plan9.armv5 \
	dukkha.plan9.armv6

dukkha.solaris.amd64:
	sh scripts/build/build.sh $@

dukkha.aix.ppc64:
	sh scripts/build/build.sh $@

dukkha.dragonfly.amd64:
	sh scripts/build/build.sh $@
