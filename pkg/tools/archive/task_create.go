package archive

import (
	"fmt"
	"io"
	"os"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindCreate = "create"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindCreate,
		func(toolName string) dukkha.Task {
			t := &TaskCreate{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindCreate, t)
			return t
		},
	)
}

type TaskCreate struct {
	rs.BaseField `yaml:"-"`

	tools.BaseTask `yaml:",inline"`

	// Format of the archive, one of [tar, zip]
	//
	// Defaults to `"zip"` when matrix.kernel is set to windows, otherwise `"tar"`
	Format string `yaml:"format"`

	// Compression configuration
	Compression *compressionSpec `yaml:"compression"`

	// Output archive file
	Output string `yaml:"output"`

	// Files to include into archive
	Files []*archiveFileSpec `yaml:"files"`
}

type compressionSpec struct {
	rs.BaseField

	// Enable compression
	//
	// Defaults to `false`
	Enabled bool `yaml:"enabled"`

	// Method of compression
	//
	// for `tar`, one of [gzip, bzip2, zstd, lzma, xz, zstd]
	// for `zip`, one of [deflate, bzip2, zstd, lzma, xz, zstd]
	//
	// Defaults to `"defalte"` when format is zip
	// Defaults to `"gzip"` when format is tar
	Method string `yaml:"method"`

	// Level of compression, value is method dependent, usually 1 - 9
	//
	// Defaults to `5`
	Level string `yaml:"level"`
}

type archiveFileSpec struct {
	rs.BaseField

	// From local file path, include files to be archived with glob pattern support
	From string `yaml:"from"`

	// To is the in archive path those files will go
	To string `yaml:"to"`
}

func (c *TaskCreate) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec

	err := c.DoAfterFieldsResolved(rc, -1, true, func() error {
		output := c.Output
		files, err := collectFiles(c.Files)
		if err != nil {
			return err
		}

		var (
			method *string
			level  string

			format = c.Format
		)

		if len(format) == 0 {
			switch rc.MatrixKernel() {
			case constant.KERNEL_WINDOWS:
				format = "zip"
			default:
				format = "tar"
			}
		}

		if c.Compression != nil && c.Compression.Enabled {
			cmethod := c.Compression.Method
			level = c.Compression.Level
			if len(cmethod) == 0 {
				switch format {
				case "zip":
					cmethod = constant.ZipCompressionMethod_Deflate.String()
					level = "5"
				case "tar":
					cmethod = "gzip"
					level = "5"
				}
			}

			method = &cmethod
		}

		steps = append(steps, dukkha.TaskExecSpec{
			AlterExecFunc: func(
				replace dukkha.ReplaceEntries,
				stdin io.Reader,
				stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				out, err := os.OpenFile(output, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return nil, err
				}

				defer func() { _ = out.Close() }()

				switch format {
				case "tar":
					var tarball io.WriteCloser = out

					if method != nil {
						tarball, err = createCompressionStream(out, *method, level)
						if err != nil {
							return nil, err
						}
					}

					err = createTar(tarball, files)
					_ = tarball.Close()

					return nil, err
				case "zip":
					return nil, createZip(out, files, method, level)
				default:
					return nil, fmt.Errorf("unsupported format: %q", format)
				}
			},
		})

		return nil
	})

	return steps, err
}
