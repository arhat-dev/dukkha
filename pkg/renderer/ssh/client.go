package ssh

import (
	"crypto/rand"
	"fmt"
	"net"
	"strconv"

	"arhat.dev/rs"
	"golang.org/x/crypto/ssh"
)

type Spec struct {
	rs.BaseField `yaml:"-"`

	// User for git ssh service, defaults to `git`
	User string `yaml:"user"`

	// Host of git ssh server e.g. gitlab.com
	Host string `yaml:"host"`

	// Port of ssh service, defaults to `22`
	Port int `yaml:"port"`

	// HostKey is the public key to verify remote host
	HostKey string `yaml:"host_key"`

	// authentication
	PrivateKey string `yaml:"private_key"`
	Password   string `yaml:"password"`
}

func (s *Spec) Clone() *Spec {
	return &Spec{
		User:       s.User,
		Host:       s.Host,
		Port:       s.Port,
		HostKey:    s.Host,
		PrivateKey: s.PrivateKey,
		Password:   s.Password,
	}
}

func NewClient(f *Spec) (*ssh.Client, error) {
	var authMethod ssh.AuthMethod

	switch {
	case len(f.PrivateKey) != 0:
		signer, err := ssh.ParsePrivateKey([]byte(f.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %w", err)
		}

		authMethod = ssh.PublicKeys(signer)
	case len(f.Password) != 0:
		authMethod = ssh.Password(f.Password)
	default:
		// TBD: allow user interaction or just return error?
		authMethod = ssh.KeyboardInteractive(
			func(
				user, instruction string,
				questions []string,
				echos []bool,
			) (answers []string, err error) {
				return
			},
		)
		_ = authMethod

		return nil, fmt.Errorf("no password or private key provided")
	}

	var hostKeyCallback ssh.HostKeyCallback
	if len(f.HostKey) != 0 {
		hostKey, err := ssh.ParsePublicKey([]byte(f.HostKey))
		if err != nil {
			return nil, fmt.Errorf("invalid host key: %w", err)
		}

		hostKeyCallback = ssh.FixedHostKey(hostKey)
	} else {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	user := f.User
	if len(user) == 0 {
		user = "git"
	}

	port := "22"
	if f.Port != 0 {
		port = strconv.FormatInt(int64(f.Port), 10)
	}

	client, err := ssh.Dial(
		"tcp",
		net.JoinHostPort(f.Host, port),
		&ssh.ClientConfig{
			Config: ssh.Config{
				Rand: rand.Reader,
			},
			Auth: []ssh.AuthMethod{authMethod},
			BannerCallback: func(message string) error {
				// TBD: what to do?
				return nil
			},
			User:            user,
			HostKeyCallback: hostKeyCallback,
			// fake version
			ClientVersion: "SSH-2.0-OpenSSH_8.1",
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to dial remote git ssh server: %w",
			err,
		)
	}

	return client, nil
}
