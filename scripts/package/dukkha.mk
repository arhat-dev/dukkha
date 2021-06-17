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

#
# linux
#
package.dukkha.deb.amd64:
	sh scripts/package/package.sh $@

package.dukkha.deb.armv6:
	sh scripts/package/package.sh $@

package.dukkha.deb.armv7:
	sh scripts/package/package.sh $@

package.dukkha.deb.arm64:
	sh scripts/package/package.sh $@

package.dukkha.deb.all: \
	package.dukkha.deb.amd64 \
	package.dukkha.deb.armv6 \
	package.dukkha.deb.armv7 \
	package.dukkha.deb.arm64

package.dukkha.rpm.amd64:
	sh scripts/package/package.sh $@

package.dukkha.rpm.armv7:
	sh scripts/package/package.sh $@

package.dukkha.rpm.arm64:
	sh scripts/package/package.sh $@

package.dukkha.rpm.all: \
	package.dukkha.rpm.amd64 \
	package.dukkha.rpm.armv7 \
	package.dukkha.rpm.arm64

package.dukkha.linux.all: \
	package.dukkha.deb.all \
	package.dukkha.rpm.all

#
# windows
#

package.dukkha.msi.amd64:
	sh scripts/package/package.sh $@

package.dukkha.msi.arm64:
	sh scripts/package/package.sh $@

package.dukkha.msi.all: \
	package.dukkha.msi.amd64 \
	package.dukkha.msi.arm64

package.dukkha.windows.all: \
	package.dukkha.msi.all

#
# darwin
#

package.dukkha.pkg.amd64:
	sh scripts/package/package.sh $@

package.dukkha.pkg.arm64:
	sh scripts/package/package.sh $@

package.dukkha.pkg.all: \
	package.dukkha.pkg.amd64 \
	package.dukkha.pkg.arm64

package.dukkha.darwin.all: \
	package.dukkha.pkg.all
