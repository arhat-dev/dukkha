package helm

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"arhat.dev/pkg/textquery"
	"github.com/itchyny/gojq"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindIndex = "index"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^helm(:.+){0,1}:index$`),
		func(subMatches []string) interface{} {
			t := &TaskIndex{}
			if len(subMatches) != 0 {
				t.SetToolName(strings.TrimPrefix(subMatches[0], ":"))
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskIndex)(nil)

type TaskIndex struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	RepoURL     string `yaml:"repo_url"`
	PackagesDir string `yaml:"packages_dir"`
	Output      string `yaml:"output"`

	PackageBaseURL string `yaml:"package_base_url"`
}

func (c *TaskIndex) ToolKind() string { return ToolKind }
func (c *TaskIndex) TaskKind() string { return TaskKindIndex }

func (c *TaskIndex) GetExecSpecs(ctx *field.RenderingContext, helmCmd []string) ([]tools.TaskExecSpec, error) {
	cacheDir := ctx.Values().Env[constant.ENV_DUKKHA_CACHE_DIR]
	indexDir, err := ioutil.TempDir(cacheDir, "helm-index-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary index dir: %w", err)
	}

	indexCmd := sliceutils.NewStrings(helmCmd, "repo", "index")

	if len(c.RepoURL) != 0 {
		indexCmd = append(indexCmd, "--url", c.RepoURL)
	}

	dukkhaWorkingDir := ctx.Values().Env[constant.ENV_DUKKHA_WORKING_DIR]
	if len(c.PackagesDir) != 0 {
		pkgDir, err2 := filepath.Abs(c.PackagesDir)
		if err2 != nil {
			return nil, fmt.Errorf("failed to determine absolute path of package_dir: %w", err2)
		}

		indexCmd = append(indexCmd, pkgDir)
	} else {
		indexCmd = append(indexCmd, dukkhaWorkingDir)
	}

	output := c.Output
	if len(output) == 0 {
		output = filepath.Join(dukkhaWorkingDir, "index.yaml")
	}

	f, err := os.Stat(output)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var oldIndexData []byte
	if err == nil {
		if f.IsDir() {
			return nil, fmt.Errorf("unexpected output destination is a directory")
		}

		// is file, prepare to merge
		oldIndexData, err = os.ReadFile(output)
		if err != nil {
			return nil, fmt.Errorf("failed to read old index file: %w", err)
		}
	}

	var steps []tools.TaskExecSpec

	steps = append(steps, tools.TaskExecSpec{
		Chdir:   indexDir,
		Command: indexCmd,
	})

	baseURL := c.PackageBaseURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	steps = append(steps, tools.TaskExecSpec{
		AlterExecFunc: func(
			replace map[string][]byte,
			stdin io.Reader, stdout, stderr io.Writer,
		) ([]tools.TaskExecSpec, error) {
			indexFile := filepath.Join(indexDir, "index.yaml")
			indexData, err := os.ReadFile(indexFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read index data")
			}

			// TODO: update index fields

			indexData, err = addPrefixToPackageURLs(indexData, baseURL)
			if err != nil {
				return nil, fmt.Errorf("failed to set base url: %w", err)
			}

			object := make(map[string]interface{})
			err = yaml.Unmarshal(indexData, &object)
			if err != nil {
				return nil, fmt.Errorf("failed to get processed index data: %w", err)
			}

			oldObject := make(map[string]interface{})
			err = yaml.Unmarshal(oldIndexData, &oldObject)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal old index data: %w", err)
			}

			result, err := sortPackagesByVersion(mergeMaps(oldObject, object))
			if err != nil {
				return nil, fmt.Errorf("failed to sort chart packages: %w", err)
			}

			return nil, os.WriteFile(output, result, f.Mode().Perm())
		},
	})

	return steps, nil
}

func addPrefixToPackageURLs(indexData []byte, prefix string) ([]byte, error) {
	addURLPrefixQuery := fmt.Sprintf(`.entries[] |= map(.urls[] |= "%s\(.)")`, prefix)

	result, err := textquery.YQBytes(addURLPrefixQuery, indexData)
	if err != nil {
		return nil, fmt.Errorf("failed to run prefix adding query over index data: %w", err)
	}

	return []byte(result), nil
}

func sortPackagesByVersion(indexObject map[string]interface{}) ([]byte, error) {
	sortQuery, err := gojq.Parse(`.entries[] |= sort_by(.version)`)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sort query: %w", err)
	}

	result, _, err := textquery.RunQuery(sortQuery, indexObject, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to sort index data: %w", err)
	}

	return []byte(textquery.HandleQueryResult(result, yaml.Marshal)), nil
}
