//spellchecker:words testutil
package testutil

//spellchecker:words testing
import (
	"fmt"
	"testing"
)

// TestingT may be used instead of *testing.T to allow testing test code.
// During tests, [RecordingT] should be used.
type TestingT interface {
	Errorf(format string, args ...any)
	Helper()
}

var (
	_ TestingT = (*testing.T)(nil)
	_ TestingT = (*RecordingT)(nil)
)

// RecordingT implements [TestingT], recording any calls the helper and Errorf functions.
// It is not safe to be used concurrently.
type RecordingT struct {
	Message      string
	HelperCalled bool
}

func (t *RecordingT) Helper() {
	t.HelperCalled = true
}

func (t *RecordingT) Errorf(format string, args ...any) {
	t.Message = fmt.Sprintf(format, args...)
}
