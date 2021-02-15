package record

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRecord_Record(t *testing.T) {
	var r Record

	N := 1000 // number of elements to record
	M := 1000 // how many times to record each element

	var recordFalse uint64
	var recordTrue uint64

	wg := &sync.WaitGroup{}
	wg.Add(N * M)
	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			go func(i int) {
				defer wg.Done()
				// record the element and count how many times it was recorded
				if r.Record(i) == false {
					atomic.AddUint64(&recordFalse, 1)
				} else {
					atomic.AddUint64(&recordTrue, 1)
				}
			}(i)
		}
	}
	wg.Wait()

	wantFalse := uint64(N) // the first time each element is recorded, it wasn't recorded before.
	if recordFalse != wantFalse {
		t.Errorf("Record.Record() got false %d times, wanted %d times", recordFalse, wantFalse)
	}

	wantTrue := uint64((M - 1) * N) // every other time it was already recorded.
	if recordTrue != wantTrue {
		t.Errorf("Record.Record() got false %d times, wanted %d times", recordTrue, wantTrue)
	}

}

func TestRecord_Recorded(t *testing.T) {
	var r Record

	N := 1000      // number of elements to record
	NN := N + 1000 // number of elements to not record

	for i := 0; i < N; i++ {
		r.Record(i)
	}

	for i := 0; i < NN; i++ {
		got := r.Recorded(i)
		want := i < N
		if got != want {
			t.Errorf("Record.Recorded(%d) = %v, want = %v", i, got, want)
		}
	}

	r.Reset()

	for i := 0; i < NN; i++ {
		got := r.Recorded(i)
		want := false
		if got != want {
			t.Errorf("Record.Rest().Recorded(%d) = %v, want = %v", i, got, want)
		}
	}

}

func ExampleRecord() {
	r := Record{}

	// record a value using the Record() method
	fmt.Println(r.Record("first"))
	fmt.Println(r.Record("first"))

	// check if a value has been recorded using .Recorded()
	fmt.Println(r.Recorded("second"))
	r.Record("second")
	fmt.Println(r.Recorded("second"))

	// Reset all recorded values using .Reset()
	r.Reset()
	fmt.Println(r.Recorded("first"))
	fmt.Println(r.Recorded("second"))

	// Output:
	// false
	// true
	// false
	// true
	// false
	// false
}
