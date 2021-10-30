//go:build add_renderer_s3 || docs
// +build add_renderer_s3 docs

package addon

import (
	_ "arhat.dev/dukkha/pkg/renderer/s3"
)
