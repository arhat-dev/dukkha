package git

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"arhat.dev/dukkha/pkg/renderer/ssh"
	"arhat.dev/pkg/iohelper"
	"arhat.dev/rs"

	gossh "golang.org/x/crypto/ssh"
)

type FetchSpec struct {
	rs.BaseField

	// Repo you want to fetch from
	Repo string `yaml:"repo"`

	// Ref is the git object reference, usually branch/tag name, defaults to `HEAD`
	Ref string `yaml:"ref"`

	// Path of the target file
	Path string `yaml:"path"`
}

type inputFetchSpec struct {
	rs.BaseField

	ssh.Spec  `yaml:",inline"`
	FetchSpec `yaml:",inline"`
}

func (f *FetchSpec) fetchRemote(sshConfig *ssh.Spec) ([]byte, error) {
	if len(f.Path) == 0 {
		return nil, fmt.Errorf("invalid no path in repo set")
	}

	client, err := ssh.NewClient(sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create ssh client: %w", err)
	}

	defer func() { _ = client.Close() }()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to open ssh seesion: %w", err)
	}
	defer func() { _ = session.Close() }()

	var (
		stdin  io.Writer
		stdout = &gitWireReader{}

		stderr = &bytes.Buffer{}
	)

	session.Stdin, stdin = iohelper.Pipe()
	stdout.reader, session.Stdout = iohelper.Pipe()
	session.Stderr = stderr

	session.Setenv("GIT_PROTOCOL", "version=2")
	err = session.Start(fmt.Sprintf("git-upload-archive '%s'", f.Repo))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to run git-upload-archive in remote host: %w",
			err,
		)
	}

	// provide object ref, usually a branch/tag name
	ref := f.Ref
	if len(ref) == 0 {
		// use default branch if not set
		ref = "HEAD"
	}

	_, err = stdin.Write([]byte(
		formatPktLines([]string{
			// provide fake --format
			//
			// ref: https://github.com/git/git/blob/master/builtin/archive.c#L39-L49
			"argument --format=tar\n",
			"argument " + ref + "\n",
			"argument " + f.Path + "\n",
		}) + "0000",
	))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to write params for git-upload-archive: %w",
			err,
		)
	}

	ackBytes, err := stdout.ReadPkt()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to check ack from git-upload-archive: %w",
			err,
		)
	}

	ackStr := string(ackBytes)
	switch {
	case strings.HasPrefix(ackStr, "ACK"):
		// success
	case strings.HasPrefix(ackStr, "NACK "):
		// failed
		return nil, fmt.Errorf(
			"git-upload-archive error: %q",
			ackStr[5:],
		)
	default:
		return nil, fmt.Errorf(
			"unexpected non ack line %q",
			ackStr,
		)
	}

	err = stdout.ReadFlush()
	if err != nil {
		return nil, fmt.Errorf(
			"unepxected not flush packet: %w", err,
		)
	}

	// now switch to recv_sideband
	sbr, err := stdout.SideBandReader()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create sideband reader: %w",
			err,
		)
	}

	tr := tar.NewReader(sbr)
	var dataBuf []byte
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read tar header: %w", err)
		}

		_ = hdr.Name
		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to read content in tar: %w", err,
			)
		}

		if len(data) != 0 {
			dataBuf = data
		}
	}

	err = session.Wait()
	if err != nil && !errors.Is(err, &gossh.ExitMissingError{}) {
		return nil, fmt.Errorf(
			"git-upload-archive exited with error: %q",
			string(stderr.Next(stderr.Len())),
		)
	}

	return dataBuf, nil
}
