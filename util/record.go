package util

import (
	"sync"
)

// Record is an object that keeps track of recorded values.
// It is safe for concurrent access, the zero value is ready to use.
//
// A Record must not be copied after first use.
type Record struct {
	m sync.RWMutex

	records map[interface{}]struct{}
}

// Record records and marks the value v as having been visited.
//
// When value has been recorded before, returns a true.
// Otherwise returns false.
//
// Record is an atomic operation and can be safely called concurrently.
func (r *Record) Record(v interface{}) (recorded bool) {
	r.m.Lock()
	defer r.m.Unlock()

	// don't use .Recorded() because of nested locking
	_, recorded = r.records[v]
	if recorded { // fast path: already recorded
		return
	}

	if r.records == nil {
		r.records = make(map[interface{}]struct{})
	}

	r.records[v] = struct{}{}
	return
}

// Recorded checks and returns if the value v has been recorded.
func (r *Record) Recorded(v interface{}) (visited bool) {
	r.m.RLock()
	defer r.m.RUnlock()

	_, visited = r.records[v]
	return
}

// Reset clears all recorded values from this Record, resetting it to an empty state.
func (r *Record) Reset() {
	r.m.Lock()
	defer r.m.Unlock()

	r.records = make(map[interface{}]struct{})
}
