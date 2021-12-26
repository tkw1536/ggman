package sema

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/tkw1536/ggman/internal/testutil"
)

func ExampleNewSemaphore() {
	// create a new semaphore with two elements
	sema := NewSemaphore(2)

	// some very finite resource pool
	var resource uint64 = 2

	// create N = 100 workers that each attempt to use the finite resource
	N := 100
	var worked uint64
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			// accounting: keep track that we did some work and that we're done!
			defer wg.Done()
			defer atomic.AddUint64(&worked, 1)

			// Lock the semaphore
			// the lock can be locked at most twice
			sema.Lock()
			defer sema.Unlock()

			// check that the resource is available
			// since we are protected by the semaphore, this is guaranteed to be the case
			if atomic.LoadUint64(&resource) == 0 {
				panic("no resource available")
			}

			// while we are working, take the resources away
			atomic.AddUint64(&resource, ^uint64(0))
			defer atomic.AddUint64(&resource, 1)

			// ... deep computation ...
			time.Sleep(10 * time.Millisecond)
		}()
	}

	wg.Wait()

	fmt.Printf("Worked %d times", worked)
	// Output: Worked 100 times
}

func ExampleNewSemaphore_simple() {
	sema := NewSemaphore(2)

	// we can lock it two times
	sema.Lock()
	sema.Lock()

	// this call would block
	// sema.Lock()
	fmt.Println("two lock calls")

	// before need to unlock to aquire again
	sema.Unlock()
	sema.Lock()

	fmt.Println("another lock call only after unlock")

	// Output: two lock calls
	// another lock call only after unlock
}
func ExampleNewSemaphore_zero() {
	// a zero or negative limit creates a semaphore without any limits
	sema := NewSemaphore(0)

	N := 1000

	// so we can call Lock as many times as we want
	for i := 0; i < N; i++ {
		sema.Lock()
	}

	// and nothing was blocked!
	fmt.Println("nothing blocked")
	// Output: nothing blocked
}

func ExampleNewSemaphore_one() {
	// a semaphore with value >= 2 is a regular semaphore
	sema := NewSemaphore(2)
	nothing := time.Nanosecond

	// do a bunch of locks and unlocks
	N := 1000

	// can lock it twice, before requiring an unlock
	for i := 0; i < N; i++ {
		sema.Lock()
		sema.Lock()

		time.Sleep(nothing)

		sema.Unlock()
		sema.Unlock()
	}

	fmt.Println("nothing blocked")
	// Output: nothing blocked
}

func ExampleNewSemaphore_two() {
	// a semaphore with value one is just a mutex
	sema := NewSemaphore(1)
	nothing := time.Nanosecond

	fmt.Printf("type = %s\n", reflect.TypeOf(sema).String())

	// do a bunch of locks and unlocks
	N := 1000

	for i := 0; i < N; i++ {
		sema.Lock()
		time.Sleep(nothing)
		sema.Unlock()
	}

	// and nothing was blocked!
	fmt.Println("nothing blocked")
	// Output: type = *sync.Mutex
	// nothing blocked
}

func ExampleNewSemaphore_panic() {
	sema := NewSemaphore(2)

	// an unlock without a corresponding unlock will always panic
	didPanic, value := testutil.DoesPanic(func() {
		sema.Unlock()
	})
	if !didPanic {
		panic("did not panic")
	}

	fmt.Printf("Unlock() panic = %#v", value)
	// Output: Unlock() panic = "Semaphore: Unlock without Lock"
}

func TestNewSemaphore_simple(t *testing.T) {
	sema := NewSemaphore(2)
	sema.Lock()
	sema.Lock()

	go func() {
		sema.Lock()
		panic("never reached")
	}()

	time.Sleep(100 * time.Microsecond)
}

func TestNewSemaphore_exhausting(t *testing.T) {
	// this test tests all cases for 1 <= n < 100
	for n := 1; n <= 100; n++ {
		s := NewSemaphore(n)

		// fully lock it
		for i := 0; i < n; i++ {
			s.Lock()
		}

		// unlock and relock one of them
		s.Unlock()
		s.Lock()
	}
}

func BenchmarkNewSemaphore_uncontested(b *testing.B) {
	sema := NewSemaphore(2)
	nothing := time.Nanosecond

	for i := 0; i < b.N; i++ {
		sema.Lock()
		sema.Lock()

		time.Sleep(nothing)

		sema.Unlock()
		sema.Unlock()
	}
}

func BenchmarkNewSemaphore_contested(b *testing.B) {
	sema := NewSemaphore(2)
	nothing := time.Nanosecond

	sema.Lock()

	// contest the semaphore in a concurrent goroutine
	go func() {
		for i := 0; i < b.N; i++ {
			sema.Lock()
			time.Sleep(nothing)

			// time.Sleep(time.Millisecond)

			sema.Unlock()
		}
	}()

	// do the attempting to aquire
	for i := 0; i < b.N; i++ {
		sema.Lock()
		time.Sleep(nothing)
		sema.Unlock()
	}
}
