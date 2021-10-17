package git

import (
	"testing"

	"arhat.dev/dukkha/pkg/renderer/ssh"
	"github.com/stretchr/testify/assert"
)

func TestFetchSpec(t *testing.T) {
	// TODO: enable fetch test
	t.SkipNow()

	spec := &FetchSpec{
		Repo: "",
		Path: "README.md",
	}

	data, err := spec.fetchRemote(&ssh.Spec{
		User:       "git",
		PrivateKey: ``,
		Host:       "gitlab.com",
	})
	assert.NoError(t, err)

	t.Log(string(data))
}
