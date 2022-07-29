package conf

import (
	"sync"
)

type sigch = chan struct{}

type Job struct {
	prevDone, nextRun sigch
}

// NewSyncGroup creates a new SyncGroup
func NewSyncGroup() (ret *SyncGroup) {
	ret = &SyncGroup{}
	ret.Init()
	return
}

// SyncGroup is a helper to implement sequential execution of parallel jobs
//
// NOTE: Init MUST be called before calling other methods
type SyncGroup struct {
	last     sigch
	canceled sigch
	wg       sync.WaitGroup
	lastErr  error
}

// Init this sync group
func (sg *SyncGroup) Init() {
	sig := make(sigch)
	close(sig)
	sg.last = sig
	sg.canceled = make(chan struct{})
}

// NewJob creates a new job handle in this sync group
func (sg *SyncGroup) NewJob() (ret Job) {
	ret.prevDone = sg.last

	next := make(sigch)
	ret.nextRun = next
	sg.last = next
	sg.wg.Add(1)
	return
}

// Go spawns a new goroutine in the sync group, user func MUST call Lock with
// the job arg passed
//
// on error return of user func, it calls Cancel and Done, otherwise, calls Unlock and Done
func (sg *SyncGroup) Go(fn func(j Job) error) {
	go func(j Job) {
		err := fn(j)
		if err != nil {
			sg.Cancel(err)
			sg.Done()
			return
		}

		sg.Unlock(j)
		sg.Done()
	}(sg.NewJob())
}

// Done is sync.WaitGroup.Done
func (sg *SyncGroup) Done() {
	sg.wg.Done()
}

// Lock waits until this job can be continued
func (sg *SyncGroup) Lock(j Job) bool {
	select {
	case <-j.prevDone:
		return true
	case <-sg.canceled:
		return false
	}
}

// Unlock wakes next waiting goroutine
func (sg *SyncGroup) Unlock(j Job) {
	select {
	case <-sg.canceled:
	default:
		close(j.nextRun)
	}
}

// Cancel the sync group with error, wakeup all waiting goroutines
func (sg *SyncGroup) Cancel(err error) {
	select {
	case <-sg.canceled:
	default:
		sg.lastErr = err
		close(sg.canceled)
	}
}

func (sg *SyncGroup) Wait() {
	select {
	case <-sg.canceled:
	default:
		sg.wg.Wait()
	}
}

// Err returns the last error stored when Cancel called
func (sg *SyncGroup) Err() error { return sg.lastErr }
