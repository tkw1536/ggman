// Package sema implements semaphores and semaphore-related scheduling
package sema

import "sync"

// NewSemaphore creates a new semaphore with the provided limit.
//
// A semaphore is a sync.Locker which guards a finite resource (of value limit).
// The resource can be aquired using a call to Lock() and released using a call to Unlock().
//
// A semaphore with limit = 1 is a *sync.Mutex.
// Semaphores with limit <= 0 implement Lock and Unlock as no-ops.
//
// Calls to Lock() block until resources in the semaphore are available.
// Calls to Unlock() never block.
//
// Calls to Unlock() when the resource is at value limit are a programming error.
// They may cause a runtime panic or an error.
func NewSemaphore(limit int) sync.Locker {
	switch {
	case limit <= 0:
		return fakeLocker{}
	case limit == 1:
		return &sync.Mutex{}
	default:
		return semaphore(make(chan struct{}, limit))
	}
}

// semaphore can implement a semiphore of limit >= 2
//
// the underlying channel must be a buffered channel
// to aquire a resource, it is writter into the underlying buffer
// to release a resource, it is read from the buffer
type semaphore chan struct{}

func (s semaphore) Lock() {
	s <- struct{}{}
}

func (s semaphore) Unlock() {
	select {
	case <-s:
	default:
		panic("Semaphore: Unlock without Lock")
	}
}

// fakeLocker locks and unlocks
type fakeLocker struct{}

func (fakeLocker) Lock()   {}
func (fakeLocker) Unlock() {}
