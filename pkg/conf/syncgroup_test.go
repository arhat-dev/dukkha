package conf

import (
	"testing"
	"time"

	"arhat.dev/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestSyncGroup_WaitExecDone(t *testing.T) {
	var (
		sg       SyncGroup
		actual   []uint32
		expected []uint32
	)

	const N = 10

	sg.Init()
	sigStart := make(chan struct{}, 0)

	for i := uint32(0); i < N; i++ {
		expected = append(expected, i)
		go func(seq Job, i uint32) {
			<-sigStart
			time.Sleep(10 * time.Millisecond)
			assert.True(t, sg.Lock(seq))
			time.Sleep(time.Duration((N - i)) * 10 * time.Millisecond)
			actual = append(actual, i)
			sg.Unlock(seq)
			sg.Done()
		}(sg.NewJob(), i)
	}
	close(sigStart)

	sg.Wait()

	assert.Equal(t, expected, actual)
}

func TestSyncGroupNotifyList_WaitExecCancel(t *testing.T) {
	const N = 10

	sg := NewSyncGroup()
	sigStart := make(chan struct{}, 0)

	go func(seq Job) {
		<-sigStart
		sg.Lock(seq)

		sg.Cancel(testhelper.Error())
		sg.Done()
	}(sg.NewJob())

	for i := uint32(0); i < N; i++ {
		go func(seq Job) {
			<-sigStart
			time.Sleep(10 * time.Millisecond)
			assert.False(t, sg.Lock(seq))
			sg.Done()
		}(sg.NewJob())
	}
	close(sigStart)

	sg.Wait()
}

func TestSyncGroup_WaitGo(t *testing.T) {
	var s string
	sg := NewSyncGroup()
	// seq 1
	sg.Go(func(seq Job) error {
		assert.True(t, sg.Lock(seq))

		// seq 2
		seq2 := sg.NewJob()
		go func() {
			sg.Lock(seq2)
			s += "2"
			sg.Unlock(seq2)
			sg.Done()
		}()

		time.Sleep(time.Second)
		s += "1"
		return nil
	})

	sg.Wait()
	assert.Equal(t, "12", s)

	// seq 3
	sg.Go(func(seq Job) error {
		assert.True(t, sg.Lock(seq))
		return testhelper.Error()
	})

	// seq 4
	sg.Go(func(seq Job) error {
		assert.False(t, sg.Lock(seq))
		return nil
	})

	sg.Wait()
}

// func BenchmarkSyncGroupNotifyList(b *testing.B) {
// 	var actual int
//
// 	sg := NewSyncGroupNotifyList()
// 	sigStart := make(chan struct{}, 0)
// 	for i := 0; i < b.N; i++ {
// 		go func(seq uint32) {
// 			<-sigStart
// 			sg.Lock(seq)
// 			actual++
// 			sg.Unlock(seq)
// 		}(sg.NewJob())
// 	}
//
// 	// b.ResetTimer()
//
// 	close(sigStart)
// 	sg.Wait()
// 	b.StopTimer()
// 	assert.Equal(b, b.N, actual)
// }

func BenchmarkSyncGroup(b *testing.B) {
	var actual int

	sg := NewSyncGroup()
	// sigStart := make(chan struct{}, 0)
	for i := 0; i < b.N; i++ {
		go func(j Job) {
			// <-sigStart
			sg.Lock(j)
			actual++
			sg.Unlock(j)
			sg.Done()
		}(sg.NewJob())
	}

	// b.ResetTimer()

	// close(sigStart)
	sg.Wait()
	b.StopTimer()
	assert.Equal(b, b.N, actual)
}

func BenchmarkSyncChannelList(b *testing.B) {
	sigStart := make(chan struct{}, 0)
	pass := make([]chan int, b.N)
	for i := 0; i < b.N; i++ {
		pass[i] = make(chan int)
	}

	go func() {
		<-sigStart
		pass[0] <- 1
	}()

	for i := 1; i < b.N; i++ {
		go func(id int) {
			<-sigStart
			pass[id] <- (<-pass[id-1] + 1)
		}(i)
	}

	b.ResetTimer()
	close(sigStart)

	actual := <-pass[b.N-1]

	b.StopTimer()
	assert.Equal(b, b.N, actual)
}
