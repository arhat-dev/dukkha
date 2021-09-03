// Code generated by 'yaegi extract sync'. DO NOT EDIT.

//go:build go1.17
// +build go1.17

package stdlib

import (
	"reflect"
	"sync"
)

func init() {
	Symbols["sync/sync"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"NewCond": reflect.ValueOf(sync.NewCond),

		// type definitions
		"Cond":      reflect.ValueOf((*sync.Cond)(nil)),
		"Locker":    reflect.ValueOf((*sync.Locker)(nil)),
		"Map":       reflect.ValueOf((*sync.Map)(nil)),
		"Mutex":     reflect.ValueOf((*sync.Mutex)(nil)),
		"Once":      reflect.ValueOf((*sync.Once)(nil)),
		"Pool":      reflect.ValueOf((*sync.Pool)(nil)),
		"RWMutex":   reflect.ValueOf((*sync.RWMutex)(nil)),
		"WaitGroup": reflect.ValueOf((*sync.WaitGroup)(nil)),

		// interface wrapper definitions
		"_Locker": reflect.ValueOf((*_sync_Locker)(nil)),
	}
}

// _sync_Locker is an interface wrapper for Locker type
type _sync_Locker struct {
	IValue  interface{}
	WLock   func()
	WUnlock func()
}

func (W _sync_Locker) Lock()   { W.WLock() }
func (W _sync_Locker) Unlock() { W.WUnlock() }
