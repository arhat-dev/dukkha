package utils

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLazyValue_Get(t *testing.T) {
	const (
		testdata = "test"
	)

	var called int32
	lv := NewLazyValue(func() string {
		_ = atomic.AddInt32(&called, 1)

		time.Sleep(5 * time.Second)
		return testdata
	})

	startSig := make(chan struct{})

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			<-startSig

			assert.EqualValues(t, testdata, lv.Get())
		}()
	}

	close(startSig)
	wg.Wait()
	assert.EqualValues(t, 1, called)
}
