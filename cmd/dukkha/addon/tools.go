package addon

import (
	// Add default tools and tasks
	_ "arhat.dev/dukkha/pkg/tools/archive"
	_ "arhat.dev/dukkha/pkg/tools/buildah"
	_ "arhat.dev/dukkha/pkg/tools/cosign"
	_ "arhat.dev/dukkha/pkg/tools/docker"
	_ "arhat.dev/dukkha/pkg/tools/git"
	_ "arhat.dev/dukkha/pkg/tools/github"
	_ "arhat.dev/dukkha/pkg/tools/golang"
	_ "arhat.dev/dukkha/pkg/tools/helm"
	_ "arhat.dev/dukkha/pkg/tools/kubectl"
	_ "arhat.dev/dukkha/pkg/tools/workflow"
)
