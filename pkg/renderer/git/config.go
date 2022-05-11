package git

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"strings"

	"arhat.dev/pkg/iohelper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/renderer/ssh"
)

type FetchSpec struct {
	rs.BaseField `yaml:"-"`

	// Repo you want to fetch from
	Repo string `yaml:"repo"`

	// Ref is the git object reference, usually branch/tag name, defaults to `HEAD`
	Ref string `yaml:"ref"`

	// Path of the target file
	Path string `yaml:"path"`
}

type inputFetchSpec struct {
	rs.BaseField `yaml:"-"`

	Fetch FetchSpec `yaml:",inline"`
	SSH   *ssh.Spec `yaml:"ssh,omitempty"`
}

func (f *FetchSpec) fetchRemote(sshConfig *ssh.Spec) (io.ReadCloser, error) {
	if len(f.Path) == 0 {
		return nil, fmt.Errorf("invalid no path in repo set")
	}

	client, err := ssh.NewClient(sshConfig)
	if err != nil {
		return nil, fmt.Errorf("create ssh client: %w", err)
	}

	defer func() { _ = client.Close() }()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("open ssh seesion: %w", err)
	}
	defer func() { _ = session.Close() }()

	var (
		stdin  io.Writer
		stdout = &gitWireReader{}

		stderr bytes.Buffer
	)

	session.Stdin, stdin = iohelper.Pipe()
	stdout.reader, session.Stdout = iohelper.Pipe()
	session.Stderr = &stderr

	err = session.Setenv("GIT_PROTOCOL", "version=2")
	if err != nil {
		return nil, fmt.Errorf("set env GIT_PROTOCOL: %w", err)
	}

	err = session.Start(fmt.Sprintf("git-upload-archive '%s'", f.Repo))
	if err != nil {
		return nil, fmt.Errorf(
			"run git-upload-archive in remote host: %w",
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
			"writing params for git-upload-archive: %w",
			err,
		)
	}

	ackBytes, err := stdout.ReadPkt()
	if err != nil {
		return nil, fmt.Errorf(
			"check ack from git-upload-archive: %w",
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
			"creating sideband reader: %w", err,
		)
	}

	tr := tar.NewReader(sbr)
	_, _ = tr.Next()

	return iohelper.CustomReadCloser(tr, func() error {
		for {
			_, err = tr.Next()
			if err == io.EOF {
				break
			}
		}

		return session.Wait()
	}), nil
}
