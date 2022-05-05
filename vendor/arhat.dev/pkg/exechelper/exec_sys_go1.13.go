//go:build go1.13
// +build go1.13

package exechelper

import (
	"context"
	"io"
	"os"
	"syscall"

	_ "unsafe" // for go:linkname
)

type execCmd struct {
	// Path is the path of the command to run.
	//
	// This is the only field that must be set to a non-zero
	// value. If Path is relative, it is evaluated relative
	// to Dir.
	Path string

	// Args holds command line arguments, including the command as Args[0].
	// If the Args field is empty or nil, Run uses {Path}.
	//
	// In typical use, both Path and Args are set by calling Command.
	Args []string

	// Env specifies the environment of the process.
	// Each entry is of the form "key=value".
	// If Env is nil, the new process uses the current process's
	// environment.
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	// As a special case on Windows, SYSTEMROOT is always added if
	// missing and not explicitly set to the empty string.
	Env []string

	// Dir specifies the working directory of the command.
	// If Dir is the empty string, Run runs the command in the
	// calling process's current directory.
	Dir string

	// Stdin specifies the process's standard input.
	//
	// If Stdin is nil, the process reads from the null device (os.DevNull).
	//
	// If Stdin is an *os.File, the process's standard input is connected
	// directly to that file.
	//
	// Otherwise, during the execution of the command a separate
	// goroutine reads from Stdin and delivers that data to the command
	// over a pipe. In this case, Wait does not complete until the goroutine
	// stops copying, either because it has reached the end of Stdin
	// (EOF or a read error) or because writing to the pipe returned an error.
	Stdin io.Reader

	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	//
	// If either is an *os.File, the corresponding output from the process
	// is connected directly to that file.
	//
	// Otherwise, during the execution of the command a separate goroutine
	// reads from the process over a pipe and delivers that data to the
	// corresponding Writer. In this case, Wait does not complete until the
	// goroutine reaches EOF or encounters an error.
	//
	// If Stdout and Stderr are the same writer, and have a type that can
	// be compared with ==, at most one goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer

	// ExtraFiles specifies additional open files to be inherited by the
	// new process. It does not include standard input, standard output, or
	// standard error. If non-nil, entry i becomes file descriptor 3+i.
	//
	// ExtraFiles is not supported on Windows.
	ExtraFiles []*os.File

	// SysProcAttr holds optional, operating system-specific attributes.
	// Run passes it to os.StartProcess as the os.ProcAttr's Sys field.
	SysProcAttr *syscall.SysProcAttr

	// Process is the underlying process, once started.
	Process *os.Process

	// ProcessState contains information about an exited process,
	// available after a call to Wait or Run.
	ProcessState *os.ProcessState

	ctx             context.Context // nil means none
	lookPathErr     error           // LookPath error, if any.
	finished        bool            // when Wait was called
	childFiles      []*os.File
	closeAfterStart []io.Closer
	closeAfterWait  []io.Closer
	goroutine       []func() error
	errch           chan error // one send per goroutine
	waitDone        chan struct{}
}

//go:linkname exec_cmd_closeDescriptors os/exec.(*Cmd).closeDescriptors
func exec_cmd_closeDescriptors(c *execCmd, closers []io.Closer)

//go:linkname exec_cmd_envv os/exec.(*Cmd).envv
func exec_cmd_envv(c *execCmd) ([]string, error)

//go:linkname exec_cmd_stdin os/exec.(*Cmd).stdin
func exec_cmd_stdin(*execCmd) (*os.File, error)

//go:linkname exec_cmd_stdout os/exec.(*Cmd).stdout
func exec_cmd_stdout(*execCmd) (*os.File, error)

//go:linkname exec_cmd_stderr os/exec.(*Cmd).stderr
func exec_cmd_stderr(*execCmd) (*os.File, error)

//go:linkname exec_cmd_argv os/exec.(*Cmd).argv
func exec_cmd_argv(*execCmd) []string

//go:linkname exec_dedupEnv os/exec.dedupEnv
func exec_dedupEnv(env []string) []string

//go:linkname exec_addCriticalEnv os/exec.addCriticalEnv
func exec_addCriticalEnv(env []string) []string
