package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPopulateGlobalEnv(t *testing.T) {
	var envNames []string
	for _, e := range os.Environ() {
		envNames = append(envNames, strings.SplitN(e, "=", 2)[0])
	}

	// DO NOT os.Clearenv(), git will not work with not environment variables

	populateGlobalEnv(context.TODO())

	requiredEnv := map[string]string{
		"GIT_BRANCH":          "",
		"GIT_COMMIT":          "",
		"GIT_TAG":             "",
		"GIT_WORKSPACE_CLEAN": "",
		"GIT_DEFAULT_BRANCH":  "master",
		"TIME_YEAR":           strconv.FormatInt(int64(time.Now().Year()), 10),
		"TIME_MONTH":          strconv.FormatInt(int64(time.Now().Month()), 10),
		"TIME_DAY":            strconv.FormatInt(int64(time.Now().Day()), 10),
		"TIME_HOUR":           strconv.FormatInt(int64(time.Now().Hour()), 10),
		"TIME_MINUTE":         strconv.FormatInt(int64(time.Now().Minute()), 10),
		"TIME_SECOND":         "",
		"HOST_OS":             "",
		"HOST_OS_VERSION":     "",
		"HOST_KERNEL":         runtime.GOOS,
		"HOST_KERNEL_VERSION": "",
		"HOST_ARCH":           "",
	}

	for _, name := range envNames {
		if _, required := requiredEnv[name]; required {
			continue
		}

		os.Unsetenv(name)
	}

	assert.Equal(t, len(requiredEnv), len(os.Environ()))

	for name, expectedValue := range requiredEnv {
		t.Run(name, func(t *testing.T) {
			val, ok := os.LookupEnv(name)
			assert.True(t, ok)

			if len(expectedValue) != 0 {
				assert.Equal(t, expectedValue, val)
			}

			t.Log(name, fmt.Sprintf("%q", val))
		})
	}
}
