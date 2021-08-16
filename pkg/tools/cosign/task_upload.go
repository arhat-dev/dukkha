package cosign

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"arhat.dev/pkg/md5helper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/templateutils"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/tools/buildah"
)

const TaskKindBuild = "upload"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindBuild,
		func(toolName string) dukkha.Task {
			t := &TaskUpload{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindBuild, t)
			return t
		},
	)
}

type TaskUpload struct {
	rs.BaseField

	tools.BaseTask `yaml:",inline"`

	UploadKind string      `yaml:"kind"`
	Files      []FileSpec  `yaml:"files"`
	Signing    signingSpec `yaml:"signing"`

	ImageNames []buildah.ImageNameSpec `yaml:"image_names"`
}

type FileSpec struct {
	rs.BaseField

	Path        string `yaml:"path"`
	ContentType string `yaml:"content_type"`
}

type signingSpec struct {
	rs.BaseField

	Enabled bool `yaml:"enabled"`

	PrivateKey         string `yaml:"private_key"`
	PrivateKeyPassword string `yaml:"private_key_password"`

	Repo string `yaml:"repo"`

	Verify    *bool  `yaml:"verify"`
	PublicKey string `yaml:"public_key"`

	Annotations map[string]string `yaml:"annotations"`
}

func (s *signingSpec) genSignAndVerifySpec(
	keyFile string,
	imageName string,
	options dukkha.TaskMatrixExecOptions,
) []dukkha.TaskExecSpec {
	if !s.Enabled {
		return nil
	}

	var steps []dukkha.TaskExecSpec

	annotations := sliceutils.FormatStringMap(s.Annotations, "=", false)

	// sign
	{
		var passwordStdin io.Reader
		if len(s.PrivateKeyPassword) != 0 {
			passwordStdin = strings.NewReader(s.PrivateKeyPassword)
		}

		signCmd := sliceutils.NewStrings(
			options.ToolCmd(), "sign", "-key", keyFile, "-slot", "signature",
		)

		for _, a := range annotations {
			signCmd = append(signCmd, "-a", a)
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
			UseShell:   options.UseShell(),
			ShellName:  options.ShellName(),
		})
	}

	if s.Verify != nil && !*s.Verify {
		return steps
	}

	// ensure public key file exists
	pubKeyFile := keyFile + ".pub"
	if len(s.PublicKey) == 0 {
		// need to derive public key from the private key

		var passwordStdin io.Reader
		if len(s.PrivateKeyPassword) != 0 {
			passwordStdin = strings.NewReader(s.PrivateKeyPassword)
		}

		steps = append(steps, dukkha.TaskExecSpec{
			Stdin: passwordStdin,
			Command: sliceutils.NewStrings(
				options.ToolCmd(), "public-key", "-key", keyFile,
				"-outfile", pubKeyFile,
			),
			UseShell:  options.UseShell(),
			ShellName: options.ShellName(),
		})
	} else {
		pubKey := s.PublicKey
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

	verifyCmd := sliceutils.NewStrings(
		options.ToolCmd(), "verify", "-key", pubKeyFile, "-slot", "signature",
	)

	for _, a := range annotations {
		verifyCmd = append(verifyCmd, "-a", a)
	}

	verifyCmd = append(verifyCmd, imageName)
	steps = append(steps, dukkha.TaskExecSpec{
		Command:   verifyCmd,
		UseShell:  options.UseShell(),
		ShellName: options.ShellName(),
	})

	return steps
}

func (s *signingSpec) ensurePrivateKey(dukkhaCacheDir string) (string, error) {
	if !s.Enabled {
		return "", nil
	}

	if len(s.PrivateKey) == 0 {
		return "", fmt.Errorf("no private key provided for signing")
	}

	dir := filepath.Join(dukkhaCacheDir, "cosign")

	keyFile := filepath.Join(
		dir,
		fmt.Sprintf(
			"private-key-%s",
			hex.EncodeToString(
				md5helper.Sum([]byte(s.PrivateKey)),
			),
		),
	)

	_, err := os.Stat(keyFile)
	if err == nil {
		return keyFile, nil
	}

	if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to check cosign private_key: %w", err)
	}

	err = os.MkdirAll(dir, 0750)
	if err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("failed to ensure cosign dir: %w", err)
	}

	err = os.WriteFile(keyFile, []byte(s.PrivateKey), 0400)
	if err != nil {
		return "", fmt.Errorf("failed to save private key to temporary file: %w", err)
	}

	return keyFile, nil
}

func (c *TaskUpload) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		// ret, err := c.createExecSpecs(rc, options)
		// steps = ret
		// return err

		keyFile, err := c.Signing.ensurePrivateKey(rc.CacheDir())
		if err != nil {
			return fmt.Errorf("failed to ensure private key: %w", err)
		}

		// cosign
		kind := c.UploadKind
		if len(kind) == 0 {
			kind = "blob"
		}

		ociOS, ok := constant.GetOciOS(rc.MatrixKernel())
		if !ok {
			ociOS = rc.MatrixKernel()
		}

		ociArch, ok := constant.GetOciArch(rc.MatrixArch())
		if !ok {
			ociArch = rc.MatrixArch()
		}

		ociVariant, _ := constant.GetOciArchVariant(rc.MatrixArch())

		var ociPlatformParts []string
		if len(ociOS) != 0 {
			ociPlatformParts = append(ociPlatformParts, ociOS)
		}

		if len(ociArch) != 0 {
			ociPlatformParts = append(ociPlatformParts, ociArch)
		}

		if len(ociVariant) != 0 {
			ociPlatformParts = append(ociPlatformParts, ociVariant)
		}

		ociPlatform := strings.Join(ociPlatformParts, "/")

		uploadCmd := sliceutils.NewStrings(options.ToolCmd(), "upload", kind)
		for _, fSpec := range c.Files {
			path := fSpec.Path

			if kind != "blob" {
				uploadCmd = append(uploadCmd, "-f", path)
				continue
			}

			if len(ociPlatform) != 0 {
				path += ":" + ociPlatform
			}

			uploadCmd = append(uploadCmd, "-f", path)

			if len(fSpec.ContentType) != 0 {
				uploadCmd = append(uploadCmd, "-ct", fSpec.ContentType)
			}
		}

		for _, spec := range c.ImageNames {
			if len(spec.Image) == 0 {
				continue
			}

			imageName := templateutils.SetDefaultImageTagIfNoTagSet(
				rc, spec.Image, true,
			)

			steps = append(steps, dukkha.TaskExecSpec{
				Command:   sliceutils.NewStrings(uploadCmd, imageName),
				UseShell:  options.UseShell(),
				ShellName: options.ShellName(),
			})

			steps = append(steps,
				c.Signing.genSignAndVerifySpec(
					keyFile, imageName, options,
				)...,
			)
		}

		return nil
	})

	return steps, err
}
