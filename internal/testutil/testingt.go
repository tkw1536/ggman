//spellchecker:words testutil
package testutil

//spellchecker:words testing
import (
	"fmt"
	"testing"
)

// TestingT may be used instead of *testing.T.
// We use an interface here to allow for testing the test code.
type TestingT interface {
	Errorf(format string, args ...any)
	Helper()
}

// check that both types implement TestingT.
var _ TestingT = (*testing.T)(nil)
var _ TestingT = (*RecordingT)(nil)

// RecordingT records a message passed to Errorf() and if the helper function has been called.
type RecordingT struct {
	Message      string
	HelperCalled bool
}

// TODO: consider tests for this.

func (t *RecordingT) Helper() {
	t.HelperCalled = true
}

func (t *RecordingT) Errorf(format string, args ...any) {
	t.Message = fmt.Sprintf(format, args...)
}
