package cosign

import (
	"fmt"
	"io"
	"os"
	"strings"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/templateutils"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/tools/buildah"
)

const TaskKindSignImage = "sign-image"

func init() {
	dukkha.RegisterTask(ToolKind, TaskKindSignImage, newTaskSignImage)
}

func newTaskSignImage(toolName string) dukkha.Task {
	t := &TaskSignImage{}
	t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindSignImage, t)
	return t
}

type TaskSignImage struct {
	rs.BaseField

	tools.BaseTask `yaml:",inline"`

	Options imageSigningOptions `yaml:",inline"`

	// ImageNames
	ImageNames []buildah.ImageNameSpec `yaml:"image_names"`
}

func (c *TaskSignImage) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var ret []dukkha.TaskExecSpec
	err := c.DoAfterFieldsResolved(rc, -1, true, func() error {
		keyFile, err := c.Options.Options.ensurePrivateKey(rc.CacheDir())
		if err != nil {
			return fmt.Errorf("failed to ensure private key: %w", err)
		}

		for _, spec := range c.ImageNames {
			if len(spec.Image) == 0 {
				continue
			}

			imageName := templateutils.SetDefaultImageTagIfNoTagSet(
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
	Annotations map[string]string `yaml:"annotations"`
}

func (s *imageSigningOptions) genSignAndVerifySpec(
	keyFile string,
	imageName string,
) []dukkha.TaskExecSpec {
	var steps []dukkha.TaskExecSpec

	annotations := sliceutils.FormatStringMap(s.Annotations, "=", false)
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

		for _, a := range annotations {
			signCmd = append(signCmd, "--annotations", a)
		}

		signCmd = append(signCmd, imageName)

		var env dukkha.Env
		if len(s.Repo) != 0 {
			env = append(env, &dukkha.EnvEntry{
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
					return nil, fmt.Errorf("failed to save public file: %w", err)
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
