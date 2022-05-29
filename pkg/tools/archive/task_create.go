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
	dukkha.RegisterTask(ToolKind, TaskKindCreate, newCreateTask)
}

func newCreateTask(toolName string) dukkha.Task {
	t := &TaskCreate{}
	t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), t)
	return t
}

type TaskCreate struct {
	rs.BaseField `yaml:"-"`

	TaskName string `yaml:"name"`

	tools.BaseTask `yaml:",inline"`

	// Format of the archive, one of [tar, zip]
	//
	// Defaults to `"zip"` when matrix.kernel is set to windows, otherwise `"tar"`
	Format string `yaml:"format"`

	// Compression configuration
	Compression compressionSpec `yaml:"compression"`

	// Output archive file
	Output string `yaml:"output"`

	// Files to be archived
	Files []*fileFromToSpec `yaml:"files"`
}

type compressionSpec struct {
	rs.BaseField `yaml:"-"`

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

type fileFromToSpec struct {
	rs.BaseField `yaml:"-"`

	// From local file path, include files to be archived with glob pattern support
	From string `yaml:"from"`

	// To in archive path, files by `From` will be put at
	//
	// if multiple files was selected by `From`, `To` MUST be a dir, thus a trailing slash is REQUIRED
	//
	// if only one file was selected by `From`, when `To` ends with a slash, it's a dir
	// matched file will be put into it, otherwise, its type is determined by matched file
	To string `yaml:"to"`

	// FollowSymlink eval symlink to store actual file instead of creating symlink in archive
	FollowSymlink bool `yaml:"follow_symlink"`
}

func (c *TaskCreate) Kind() dukkha.TaskKind { return TaskKindCreate }
func (c *TaskCreate) Name() dukkha.TaskName { return dukkha.TaskName(c.TaskName) }
func (c *TaskCreate) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: c.Kind(), Name: c.Name()}
}

func (c *TaskCreate) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec

	err := c.DoAfterFieldsResolved(rc, -1, true, func() error {
		output := c.Output
		files, err := collectFiles(rc.FS(), c.Files)
		if err != nil {
			return err
		}

		var (
			format = c.Format

			enableCompression = c.Compression.Enabled
			compressionMethod string
			compressionLevel  string
		)

		if len(format) == 0 {
			switch rc.MatrixKernel() {
			case constant.KERNEL_Windows:
				format = constant.ArchiveFormat_Zip
			default:
				format = constant.ArchiveFormat_Tar
			}
		}

		if enableCompression {
			compressionMethod = c.Compression.Method
			compressionLevel = c.Compression.Level

			if len(compressionMethod) == 0 {
				switch format {
				case constant.ArchiveFormat_Zip:
					compressionMethod = constant.ZipCompressionMethod_Deflate.String()
				case constant.ArchiveFormat_Tar:
					compressionMethod = constant.CompressionMethod_Gzip
				}
			}
		}

		steps = append(steps, dukkha.TaskExecSpec{
			AlterExecFunc: func(
				replace dukkha.ReplaceEntries,
				stdin io.Reader,
				stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				archiveFile, err := rc.FS().OpenFile(output, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return nil, err
				}

				out := archiveFile.(io.WriteCloser)
				defer func() { _ = out.Close() }()

				switch format {
				case constant.ArchiveFormat_Tar:
					if enableCompression {
						out, err = createCompressionStream(out, compressionMethod, compressionLevel)
						if err != nil {
							return nil, err
						}
					}

					return nil, createTar(rc.FS(), out, files)
				case constant.ArchiveFormat_Zip:
					return nil, createZip(rc.FS(), out, files, enableCompression, compressionMethod, compressionLevel)
				default:
					return nil, fmt.Errorf("unsupported format: %q", format)
				}
			},
		})

		return nil
	})

	return steps, err
}
