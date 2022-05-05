package dukkha

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	_ctx := NewConfigResolvingContext(context.Background(), nil, nil)

	ctx := _ctx.WithCustomParent(&canceledContext{})
	select {
	case <-ctx.Done():
		assert.Error(t, ctx.Err())
	default:
		assert.Fail(t, "context parent not updated")
	}
}
