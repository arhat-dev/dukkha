package synchain

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type sigch = chan struct{}

// Ticket is the handle used to ensure sequence
type Ticket struct {
	prevDone, nextRun sigch
}

// NewSynchain creates a new Synchain
func NewSynchain() (ret *Synchain) {
	ret = &Synchain{}
	ret.Init()
	return
}

// Synchain is a helper to implement internal sequential execution of parallel jobs
//
// NOTE: Init MUST be called before calling other methods
type Synchain struct {
	last     sigch
	canceled sigch
	wg       sync.WaitGroup

	spin    uint32
	lastErr error
}

// Init this sync group
func (sg *Synchain) Init() {
	sig := make(sigch)
	close(sig)
	sg.last = sig
	sg.canceled = make(chan struct{})
}

// NewTicket creates a new sequence handle in this sync group
func (sg *Synchain) NewTicket() (ret Ticket) {
	next := make(sigch)

	for !atomic.CompareAndSwapUint32(&sg.spin, 0, 1) {
		runtime.Gosched()
	}

	ret.prevDone = sg.last
	sg.last = next

	atomic.StoreUint32(&sg.spin, 0)

	ret.nextRun = next
	sg.wg.Add(1)
	return
}

// Go spawns a new goroutine in the sync group, user func MUST call Lock with
// the ticket passed
//
// on error return of user func, it calls Cancel and Done, otherwise, calls Unlock and Done
func (sg *Synchain) Go(fn func(t Ticket) error) {
	go func(j Ticket) {
		err := fn(j)
		if err != nil {
			sg.Cancel(err)
			sg.Done()
			return
		}

		sg.Unlock(j)
		sg.Done()
	}(sg.NewTicket())
}

// Done is sync.WaitGroup.Done
func (sg *Synchain) Done() {
	sg.wg.Done()
}

// Lock waits until this job can be continued
func (sg *Synchain) Lock(t Ticket) bool {
	select {
	case <-t.prevDone:
		return true
	case <-sg.canceled:
		return false
	}
}

// Unlock wakes next waiting goroutine
func (sg *Synchain) Unlock(t Ticket) {
	select {
	case <-sg.canceled:
	default:
		close(t.nextRun)
	}
}

// Cancel the sync group with error, wakeup all waiting goroutines
func (sg *Synchain) Cancel(err error) {
	select {
	case <-sg.canceled:
	default:
		sg.lastErr = err
		close(sg.canceled)
	}
}

func (sg *Synchain) Wait() {
	select {
	case <-sg.canceled:
	default:
		sg.wg.Wait()
	}
}

// Err returns the last error stored when Cancel called
func (sg *Synchain) Err() error { return sg.lastErr }
