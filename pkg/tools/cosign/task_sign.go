package cosign

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/md5helper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindSign = "sign"

func init() {
	dukkha.RegisterTask(ToolKind, TaskKindSign, newTaskSign)
}

func newTaskSign(toolName string) dukkha.Task {
	t := &TaskSign{}
	t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), t)
	return t
}

// TaskSign signs blob
type TaskSign struct {
	rs.BaseField `yaml:"-"`

	TaskName string `yaml:"name"`

	tools.BaseTask `yaml:",inline"`

	Options blobSigningOptions `yaml:",inline"`

	// Files to sign
	Files []*blobSigningFileSpec `yaml:"files"`
}

type blobSigningFileSpec struct {
	rs.BaseField

	// Path is the local file path to the blob
	Path string `yaml:"path"`

	// Output is the destination path of signature output
	Output string `yaml:"output"`
}

func (c *TaskSign) Kind() dukkha.TaskKind { return TaskKindSign }
func (c *TaskSign) Name() dukkha.TaskName { return dukkha.TaskName(c.TaskName) }
func (c *TaskSign) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: c.Kind(), Name: c.Name()}
}

func (c *TaskSign) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var ret []dukkha.TaskExecSpec
	err := c.DoAfterFieldsResolved(rc, -1, true, func() error {
		keyFile, err := c.Options.ensurePrivateKey(c.CacheFS)
		if err != nil {
			return fmt.Errorf("failed to ensure private key: %w", err)
		}

		for _, fSpec := range c.Files {
			ret = append(ret,
				c.Options.genSignAndVerifySpec(
					keyFile, fSpec.Path, fSpec.Output,
				)...,
			)
		}

		return nil
	})

	return ret, err
}

type blobSigningOptions struct {
	// PrivateKey is the content of private key to sign content
	PrivateKey string `yaml:"private_key"`

	// PrivateKeyPassword is the password to the private key
	PrivateKeyPassword string `yaml:"private_key_password"`

	// Verify signature of signed content
	//
	// Defaults to `true`
	Verify *bool `yaml:"verify"`

	// PublicKey is the content of public key to verify signed content
	//
	// if not set, derive from private key
	PublicKey string `yaml:"public_key"`
}

func (s *blobSigningOptions) ensurePrivateKey(cacheFS *fshelper.OSFS) (string, error) {
	if len(s.PrivateKey) == 0 {
		return "", fmt.Errorf("no private key provided for signing")
	}

	keyFile := "private-key-" + hex.EncodeToString(
		md5helper.Sum([]byte(s.PrivateKey)),
	)

	_, err := cacheFS.Stat(keyFile)
	if err == nil {
		return cacheFS.Abs(keyFile)
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("failed to check cosign private_key: %w", err)
	}

	err = cacheFS.WriteFile(keyFile, []byte(s.PrivateKey), 0400)
	if err != nil {
		return "", fmt.Errorf("failed to save private key to temporary file: %w", err)
	}

	return cacheFS.Abs(keyFile)
}

func (s *blobSigningOptions) genSignAndVerifySpec(
	keyFile string,
	file string,
	signatureFile string,
) []dukkha.TaskExecSpec {
	var steps []dukkha.TaskExecSpec

	if len(signatureFile) == 0 {
		signatureFile = file + ".sig"
	}

	// sign
	{
		var passwordStdin io.Reader
		if len(s.PrivateKeyPassword) != 0 {
			passwordStdin = strings.NewReader(s.PrivateKeyPassword)
		}

		signBlobCmd := []string{
			constant.DUKKHA_TOOL_CMD,
			"sign-blob",
			"--key", keyFile,
			"--slot", "signature",
			"--output", signatureFile,
			file,
		}

		steps = append(steps, dukkha.TaskExecSpec{
			Stdin:   passwordStdin,
			Command: signBlobCmd,
		})
	}

	if s.Verify != nil && !*s.Verify {
		// verification disabled manually
		return steps
	}

	// verify signature

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
			Command: []string{
				constant.DUKKHA_TOOL_CMD,
				"public-key",
				"--key", keyFile,
				"--outfile", pubKeyFile,
			},
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

	verifyCmd := []string{
		constant.DUKKHA_TOOL_CMD,
		"verify-blob",
		"--key", pubKeyFile,
		"--slot", "signature",
		"--signature", signatureFile,
		file,
	}

	steps = append(steps, dukkha.TaskExecSpec{
		Command: verifyCmd,
	})

	return steps
}
