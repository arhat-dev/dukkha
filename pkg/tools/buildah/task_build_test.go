package buildah

import (
	"testing"

	"arhat.dev/pkg/textquery"
	"github.com/stretchr/testify/assert"

	_ "embed"
)

var (
	//go:embed testdata/manifest-spec-empty.json
	manifestSpecEmpty []byte

	//go:embed testdata/manifest-spec-null.json
	manifestSpecNull []byte

	//go:embed testdata/manifest-spec-kernel-linux-missing-arch-amd64.json
	manifestSpecKernelLinuxMissingArchAMD64 []byte
)

func TestCreateManifestOsArchVariantQueryForDigest(t *testing.T) {
	tests := []struct {
		name              string
		mKernel           string
		mArch             string
		manifestSpecInput []byte
		expectErr         bool
		expected          []string
	}{
		{
			name:              "Query Empty",
			mKernel:           "linux",
			mArch:             "amd64",
			manifestSpecInput: manifestSpecEmpty,
			expected:          []string{""},
		},
		{
			name:              "Query Null",
			mKernel:           "linux",
			mArch:             "amd64",
			manifestSpecInput: manifestSpecNull,
			expectErr:         true,
		},
		{
			name:              "Query Missing AMD64",
			mKernel:           "linux",
			mArch:             "amd64",
			manifestSpecInput: manifestSpecKernelLinuxMissingArchAMD64,
			expected:          []string{""},
		},
		{
			name:              "Query linux/arm/v5",
			mKernel:           "linux",
			mArch:             "armv5",
			manifestSpecInput: manifestSpecKernelLinuxMissingArchAMD64,
			expected:          []string{"sha256:d373bfd4642d31a0d2978dd4deeec537975e09dcde175560b8f7c98cfe640159"},
		},
		{
			name:              "Query Duplicate linux/arm/v7",
			mKernel:           "linux",
			mArch:             "armv7",
			manifestSpecInput: manifestSpecKernelLinuxMissingArchAMD64,
			expected: []string{
				"sha256:0c61e926eed09cee59238c85f8f85c05297657e490addf1a6675668fff3f0727",
				"sha256:0c61e926eed09cee59238c85f8f85c05297657e490addf1a6675668fff3f0727",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			query := createManifestOsArchVariantQueryForDigest(test.mKernel, test.mArch)

			result, err := textquery.JQBytes(query, test.manifestSpecInput)
			if test.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			digests, err := parseManifestOsArchVariantQueryResult(result)
			assert.NoError(t, err)

			if !assert.Equal(t, test.expected, digests) {
				t.Log("Query:", query)
				t.Log("Manifest:", string(test.manifestSpecInput))
			}
		})
	}
}
