package cmd

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/constant"
)

func TestCreateGlobalEnv(t *testing.T) {
	t.Parallel()

	globalEnv := createGlobalEnv(context.TODO(), ".")

	now := time.Now().Local()
	zone, offset := now.Zone()
	requiredEnv := map[string]string{
		"GIT_BRANCH":         "",
		"GIT_COMMIT":         "",
		"GIT_TAG":            "",
		"GIT_WORKTREE_CLEAN": "",
		"GIT_DEFAULT_BRANCH": "master",

		"TIME_ZONE":        zone,
		"TIME_ZONE_OFFSET": strconv.FormatInt(int64(offset), 10),
		"TIME_YEAR":        strconv.FormatInt(int64(now.Year()), 10),
		"TIME_MONTH":       strconv.FormatInt(int64(now.Month()), 10),
		"TIME_DAY":         strconv.FormatInt(int64(now.Day()), 10),
		"TIME_HOUR":        strconv.FormatInt(int64(now.Hour()), 10),
		"TIME_MINUTE":      strconv.FormatInt(int64(now.Minute()), 10),
		"TIME_SECOND":      "",

		"HOST_OS":             "",
		"HOST_OS_VERSION":     "",
		"HOST_KERNEL":         runtime.GOOS,
		"HOST_KERNEL_VERSION": "",
		"HOST_ARCH":           "",
		"HOST_ARCH_SIMPLE":    "",

		"DUKKHA_WORKDIR":   "",
		"DUKKHA_CACHE_DIR": "",
	}

	for name, expectedValue := range requiredEnv {
		t.Run(name, func(t *testing.T) {
			id := constant.GetGlobalEnvIDByName(name)
			if !assert.NotEqual(t, constant.GlobalEnv(-1), id) {
				return
			}

			if len(expectedValue) != 0 {
				assert.Equal(t, expectedValue, globalEnv[id].GetLazyValue())
			}

			t.Log(name, fmt.Sprintf("%q", globalEnv[id].GetLazyValue()))
		})
	}
}
