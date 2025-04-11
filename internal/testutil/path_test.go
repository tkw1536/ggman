//spellchecker:words testutil
package testutil_test

//spellchecker:words path filepath testing github ggman internal testutil
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/ggman/internal/testutil"
)

// copied over from package proper.
var defaultVolumeName = filepath.VolumeName(os.TempDir())

func TestToOSPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want string
	}{
		{"hello/world", "hello" + string(os.PathSeparator) + "world"},
		{"", ""},
		{"./", "." + string(os.PathSeparator)},
		{"hello/../world", "hello" + string(os.PathSeparator) + ".." + string(os.PathSeparator) + "world"},
		{"/root", defaultVolumeName + string(os.PathSeparator) + "root"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := testutil.ToOSPath(tt.name); got != tt.want {
				t.Errorf("ToOSPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
