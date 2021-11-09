package testhelper

import (
	"os"
	"syscall"
)

func HijackStandardStreams(stdin, stdout, stderr *os.File, do func()) {
	if stdout != nil {
		originalStdout := os.Stdout
		originalStdoutFD := syscall.Stdout
		defer func() {
			os.Stdout = originalStdout
			syscall.Stdout = originalStdoutFD
		}()

		os.Stdout = stdout
		syscall.Stdout = getSyscallFD(stdout.Fd())
	}

	if stderr != nil {
		originalStderr := os.Stderr
		originalStderrFD := syscall.Stderr
		defer func() {
			os.Stderr = originalStderr
			syscall.Stderr = originalStderrFD
		}()

		os.Stderr = stderr
		syscall.Stderr = getSyscallFD(stderr.Fd())
	}

	if stdin != nil {
		originalStdin := os.Stdin
		originalStdinFD := syscall.Stdin
		defer func() {
			os.Stdin = originalStdin
			syscall.Stdin = originalStdinFD
		}()

		os.Stdin = stdin
		syscall.Stdin = getSyscallFD(stdin.Fd())
	}

	do()
}
