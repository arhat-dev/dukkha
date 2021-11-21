package archive

import (
	"testing"

	"arhat.dev/rs"
)

func TestTaskCreate(t *testing.T) {
	type TestCase struct {
		rs.BaseField

		Task      *TaskCreate `yaml:"task"`
		ExpectErr bool        `yaml:"expect_err"`
	}

	type ExpectedEntry struct {
		From string `yaml:"from"`
		Link string `yaml:"link"`
	}

	// 	testhelper.TestFixtures(t, "./fixtures/create",
	// 		func() interface{} { return rs.Init(&TestCase{}, nil) },
	// 		func() interface{} { return &ExpectedEntry{} },
	// 		func(t *testing.T, spec, exp interface{}) {
	//
	// 		},
	// 	)
	_, _ = &TestCase{}, &ExpectedEntry{}
}
