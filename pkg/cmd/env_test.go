package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPopulateGlobalEnv(t *testing.T) {
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
		"HOST_OS":             runtime.GOOS,
		"HOST_ARCH":           "",
	}

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
