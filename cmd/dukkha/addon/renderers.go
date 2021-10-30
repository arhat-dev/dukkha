package addon

import (
	// Add default disabled renderers
	_ "arhat.dev/dukkha/pkg/renderer/git"
	_ "arhat.dev/dukkha/pkg/renderer/http"
	_ "arhat.dev/dukkha/pkg/renderer/input"
)
