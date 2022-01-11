package sema

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

func ExampleSchedule() {
	var counter uint64

	// because we have a parallelism of 1, we run exactly in order!
	Schedule(func(i int) error {
		atomic.AddUint64(&counter, 1)
		return nil
	}, 1000, Concurrency{
		Limit: 0,
		Force: false,
	})

	fmt.Print(counter)
	// Output: 1000
}

func ExampleSchedule_order() {
	// because we have a parallelism of 1, we run exactly in order!
	Schedule(func(i int) error {
		fmt.Print(i, ";")
		return nil
	}, 4, Concurrency{
		Limit: 1,
		Force: false,
	})

	// Output: 0;1;2;3;
}

func ExampleSchedule_error() {
	err := Schedule(func(i int) error {
		// the first invocation produces an error and returns immediatly!
		if i == 0 {
			return errors.New("first function error")
		}

		// concurrently with the first invocation, we have at most one other
		// so give the first function some time to produce an error
		time.Sleep(100 * time.Millisecond)

		// the third and fourth invocations should never be called
		// since by the time the first function finishes
		// the second one has already produced an error
		if i > 1 {
			panic("never reached")
		}
		return nil
	}, 4, Concurrency{
		Limit: 2,
		Force: false,
	})
	fmt.Println(err)

	// Output: first function error
}

func ExampleSchedule_force() {
	var counter uint64

	err := Schedule(func(i int) error {
		// count the number of invocations
		atomic.AddUint64(&counter, 1)

		// the first function returns an error
		// but because of force = True, the execution is not aborted
		if i == 0 {
			return errors.New("first function error")
		}

		// ... work ...
		time.Sleep(50 * time.Millisecond)
		return nil
	}, 10, Concurrency{
		Limit: 2,
		Force: true,
	})

	fmt.Println(err)
	fmt.Println(counter, "workers called")

	// Output: first function error
	// 10 workers called
}
