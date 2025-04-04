//spellchecker:words record
package record

//spellchecker:words sync atomic testing
import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRecord_Record(t *testing.T) {
	var r Record

	const N = 1000 // number of elements to record
	const M = 1000 // how many times to record each element

	var recordFalse atomic.Uint64
	var recordTrue atomic.Uint64

	wg := &sync.WaitGroup{}
	wg.Add(N * M)
	for i := range N {
		for range M {
			go func(i int) {
				defer wg.Done()
				// record the element and count how many times it was recorded
				if r.Record(i) == false {
					recordFalse.Add(1)
				} else {
					recordTrue.Add(1)
				}
			}(i)
		}
	}
	wg.Wait()

	// the first time each element is recorded, it wasn't recorded before.
	wantFalse := uint64(N)
	if recordFalse.Load() != wantFalse {
		t.Errorf("Record.Record() got false %d times, wanted %d times", recordFalse.Load(), wantFalse)
	}

	// every other time it was already recorded.
	wantTrue := uint64((M - 1) * N)
	if recordTrue.Load() != wantTrue {
		t.Errorf("Record.Record() got false %d times, wanted %d times", recordTrue.Load(), wantTrue)
	}
}

func TestRecord_Recorded(t *testing.T) {
	var r Record

	N := 1000      // number of elements to record
	NN := N + 1000 // number of elements to not record

	for i := range N {
		r.Record(i)
	}

	for i := range NN {
		got := r.Recorded(i)
		want := i < N
		if got != want {
			t.Errorf("Record.Recorded(%d) = %v, want = %v", i, got, want)
		}
	}

	r.Reset()

	for i := range NN {
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
