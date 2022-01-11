package sema

import (
	"sync"
	"sync/atomic"
)

// Concurrency represents the amount of concurreny of an operation
type Concurrency struct {
	Limit int  // Limit inidicates the maximum number of concurrent operations. 0 or negative implies no limit.
	Force bool // Force indicates if a failed operation should still allow future operations to start
}

// Schedule schedules count instances of worker to be called.
// When count is non-positive, no calls are scheduled and immediatly returns nil.
// Each call to worker receives a unique id, from 0 up to count (exclusive count).
//
// Workers are approximatly started in order.
// Any call to worker(i) is called before worker(j) for any i < j.
// Concurrency determines the amount of concurrency that takes place for scheduling.
//
// There is no syncronisation mechanism beyond the limits themselves.
// In particular for Limit != 1, the order guarantee might be broken:
// While the invocation of worker(i) occurs before worker(j),
// no such guarantee for the first statement within those invocation is true.
//
// If an error occurs (indicated by a non-nil return from worker),
// Schedule waits for all ongoing worker calls to return.
// It may schedule further calls to worker, as determined by concurrency.Force.
// It then returns the non-nil error which triggered the error stop.
//
// If no error occurs, schedule returns nil.
func Schedule(worker func(int) error, count int, concurrency Concurrency) (err error) {
	if count <= 0 {
		return nil
	}

	sema := NewSemaphore(concurrency.Limit) // semaphore to use for the limit

	var next uint64    // id of next worker call!
	var errDone uint64 // 0 => everything is fine, 1 => no error!

	// create a wait group that waits for all the work to be done!
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()

			// aquire the semaphore
			sema.Lock()
			defer sema.Unlock()

			// check if something already broke
			// and if so, stop doing stuff!
			if !concurrency.Force && atomic.LoadUint64(&errDone) != 0 {
				return
			}

			// grab the next id to work on
			id := int(atomic.AddUint64(&next, 1) - 1)

			// do the work!
			res := worker(id)
			if res == nil {
				return
			}

			// store the error (if we haven't already)
			if atomic.CompareAndSwapUint64(&errDone, 0, 1) {
				err = res
			}
		}()
	}

	wg.Wait()
	return err
}
