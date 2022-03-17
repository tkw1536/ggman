// Package sema implements semaphores and semaphore-related scheduling
package sema

import "sync"

// NewSemaphore creates a new semaphore with the provided limit.
//
// A semaphore with limit = 1 is a *sync.Mutex.
// Semaphores with limit <= 0 implement Lock, TryLock and Unlock as no-ops that return immediatly.
func NewSemaphore(limit int) Semaphore {
	switch {
	case limit <= 0:
		return fakeLocker{}
	case limit == 1:
		return &sync.Mutex{}
	default:
		return semaphore(make(chan struct{}, limit))
	}
}

// Semaphore represents a guards of a finite resource with a specific limit.
// The resource can be acquired using a call to Lock() and released using a call to Unlock().
//
// See also NewSemaphore and sync.Locker.
type Semaphore interface {

	// Lock atomically acquires a unit of the guarded resource.
	// When the resource is not available, it blocks until such a resource is available.
	Lock()

	// TryLock attempts to atomically aquire the resource without locking.
	// When it suceeds, it returns true, otherwise it returns false.
	//
	// Calls to TryLock() never block; they always return immediatly.
	TryLock() bool

	// Unlock releases one unit of the resource that has been previously acquired.
	// Calls to Unlock() never block.
	//
	// Calls to Unlock() without an acquired resource are a programming error;
	// they may produce a panic() or a runtime error.
	Unlock()
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

func (s semaphore) TryLock() bool {
	select {
	case s <- struct{}{}:
		return true
	default:
		return false
	}
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

func (fakeLocker) Lock()         {}
func (fakeLocker) Unlock()       {}
func (fakeLocker) TryLock() bool { return true }
