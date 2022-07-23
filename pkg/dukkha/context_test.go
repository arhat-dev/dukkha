package dukkha

import (
	"context"
	"fmt"
	"testing"
	"time"

	"arhat.dev/tlang"
	"github.com/stretchr/testify/assert"
)

var (
	_ ConfigResolvingContext = (*dukkhaContext)(nil)

	_ Context = (*dukkhaContext)(nil)
)

type canceledContext struct{}

func (*canceledContext) Done() <-chan struct{} {
	ret := make(chan struct{})
	close(ret)
	return ret
}

func (*canceledContext) Deadline() (time.Time, bool)   { return time.Time{}, false }
func (*canceledContext) Err() error                    { return fmt.Errorf("canceled") }
func (*canceledContext) Value(interface{}) interface{} { return nil }

func TestContext_SetCustomParent(t *testing.T) {
	t.Parallel()

	fakeGlobalEnv := &GlobalEnvSet{}
	for i := range fakeGlobalEnv {
		fakeGlobalEnv[i] = tlang.ImmediateString("")
	}

	_ctx := NewConfigResolvingContext(context.Background(), nil, fakeGlobalEnv)

	ctx := _ctx.WithCustomParent(&canceledContext{})
	select {
	case <-ctx.Done():
		assert.Error(t, ctx.Err())
	default:
		assert.Fail(t, "context parent not updated")
	}
}
