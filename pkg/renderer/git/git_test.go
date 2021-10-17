package git

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/renderer/ssh"
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
