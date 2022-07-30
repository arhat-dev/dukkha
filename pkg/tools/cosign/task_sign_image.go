package cosign

import (
	"fmt"
	"io"
	"os"
	"strings"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/templateutils"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/tools/buildah"
)

const TaskKindSignImage = "sign-image"

func init() {
	dukkha.RegisterTask(ToolKind, TaskKindSignImage, tools.NewTask[TaskSignImage, *TaskSignImage])
}

type TaskSignImage struct {
	tools.BaseTask[CosignSignImage, *CosignSignImage]
}

// nolint:revive
type CosignSignImage struct {
	Options imageSigningOptions `yaml:",inline"`

	// ImageNames
	ImageNames []buildah.ImageNameSpec `yaml:"image_names"`

	parent tools.BaseTaskType
}

func (c *CosignSignImage) ToolKind() dukkha.ToolKind       { return ToolKind }
func (c *CosignSignImage) Kind() dukkha.TaskKind           { return TaskKindSignImage }
func (c *CosignSignImage) LinkParent(p tools.BaseTaskType) { c.parent = p }

func (c *CosignSignImage) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var ret []dukkha.TaskExecSpec
	err := c.parent.DoAfterFieldsResolved(rc, -1, true, func() error {
		keyFile, err := c.Options.Options.ensurePrivateKey(c.parent.CacheFS())
		if err != nil {
			return fmt.Errorf("ensuring private key: %w", err)
		}

		for _, spec := range c.ImageNames {
			if len(spec.Image) == 0 {
				continue
			}

			imageName := templateutils.GetFullImageName_UseDefault_IfIfNoTagSet(
				rc, spec.Image, true,
			)

			ret = append(ret,
				c.Options.genSignAndVerifySpec(
					keyFile,
					imageName,
				)...,
			)
		}

		return nil
	})

	return ret, err
}

type imageSigningOptions struct {
	rs.BaseField

	Options blobSigningOptions `yaml:",inline"`

	// Repo is the signature storage repo, defaults to the same repo as
	// image name
	Repo string `yaml:"repo"`

	// Annotations are additional key-value data pairs added when signing
	Annotations []*dukkha.NameValueEntry `yaml:"annotations"`
}

func (s *imageSigningOptions) genSignAndVerifySpec(
	keyFile string,
	imageName string,
) []dukkha.TaskExecSpec {
	var steps []dukkha.TaskExecSpec

	annotations := make([]string, len(s.Annotations))
	for i, a := range s.Annotations {
		annotations[i] = a.Name + "=" + a.Value
	}

	// sign
	{
		var passwordStdin io.Reader
		if len(s.Options.PrivateKeyPassword) != 0 {
			passwordStdin = strings.NewReader(s.Options.PrivateKeyPassword)
		}

		signCmd := []string{
			constant.DUKKHA_TOOL_CMD,
			"sign",
			"--key", keyFile,
			"--slot", "signature",
		}

		for _, anno := range annotations {
			signCmd = append(signCmd, "--annotations", anno)
		}

		signCmd = append(signCmd, imageName)

		var env dukkha.NameValueList
		if len(s.Repo) != 0 {
			env = append(env, &dukkha.NameValueEntry{
				Name:  "COSIGN_REPOSITORY",
				Value: s.Repo,
			})
		}

		steps = append(steps, dukkha.TaskExecSpec{
			EnvSuggest: env,
			Stdin:      passwordStdin,
			Command:    signCmd,
		})
	}

	if s.Options.Verify != nil && !*s.Options.Verify {
		// verification disabled manually
		return steps
	}

	// requested verification

	// ensure public key file exists
	pubKeyFile := keyFile + ".pub"
	if len(s.Options.PublicKey) == 0 {
		// need to derive public key from the private key

		var passwordStdin io.Reader
		if len(s.Options.PrivateKeyPassword) != 0 {
			passwordStdin = strings.NewReader(s.Options.PrivateKeyPassword)
		}

		steps = append(steps, dukkha.TaskExecSpec{
			Stdin: passwordStdin,
			Command: []string{
				constant.DUKKHA_TOOL_CMD,
				"public-key",
				"--key", keyFile,
				"--outfile", pubKeyFile,
			},
		})
	} else {
		pubKey := s.Options.PublicKey
		steps = append(steps, dukkha.TaskExecSpec{
			AlterExecFunc: func(
				replace dukkha.ReplaceEntries,
				stdin io.Reader,
				stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				err := os.WriteFile(pubKeyFile, []byte(pubKey), 0644)
				if err != nil {
					return nil, fmt.Errorf("saving public file: %w", err)
				}
				return nil, nil
			},
		})
	}

	verifyCmd := []string{
		constant.DUKKHA_TOOL_CMD,
		"verify",
		"--key", pubKeyFile,
		"--slot", "signature",
	}

	for _, anno := range annotations {
		verifyCmd = append(verifyCmd, "--annotations", anno)
	}

	verifyCmd = append(verifyCmd, imageName)
	steps = append(steps, dukkha.TaskExecSpec{
		Command: verifyCmd,
	})

	return steps
}
