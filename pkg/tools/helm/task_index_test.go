package helm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	_ "embed"
)

var (
	//go:embed testdata/000-test_index_data.yaml
	testIndexData string

	//go:embed testdata/100-expected_with_package_url_prefix.yaml
	expectedWithPackageURLPrefix string

	//go:embed testdata/200-expected_sorted_by_version.yaml
	expectedSortedByVersion string
)

func TestAddPrefixToPackageURLs(t *testing.T) {
	result, err := addPrefixToPackageURLs([]byte(testIndexData), "new/")
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t,
		strings.TrimSpace(expectedWithPackageURLPrefix),
		strings.TrimSpace(string(result)),
	)
}

func TestSortPackagesByVersion(t *testing.T) {
	m := make(map[string]interface{})
	yaml.Unmarshal([]byte(testIndexData), &m)

	result, err := sortPackagesByVersion(m)
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t,
		strings.TrimSpace(expectedSortedByVersion),
		strings.TrimSpace(string(result)),
	)
}
