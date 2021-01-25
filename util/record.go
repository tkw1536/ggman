package util

import (
	"sync"
)

// Record is an object that keeps track of recorded values.
// It is safe for concurrent access, the zero value is ready to use.
//
// Re-use of the Record is possible (using the Reset method), however
// the implementation is internally optimized when writes happen only once.
//
// A Record must not be copied after first use.
type Record struct {
	records sync.Map
}

// Record records and marks the value v as having been visited.
//
// When value has been recorded before, returns a true.
// Otherwise returns false.
//
// Record is an atomic operation and can be safely called concurrently.
func (r *Record) Record(v interface{}) (recorded bool) {
	_, recorded = r.records.LoadOrStore(v, struct{}{})
	return
}

// Recorded checks and returns if the value v has been recorded.
//
// Recorded is an atomic operation and can be safely called concurrently.
func (r *Record) Recorded(v interface{}) (visited bool) {
	_, visited = r.records.Load(v)
	return
}

// Reset clears all recorded values from this Record, resetting it to an empty state.
//
// Reset is not safe for concurrent usage; no reads or writes should happen concurrently with a reset.
func (r *Record) Reset() {
	r.records = sync.Map{}
}
