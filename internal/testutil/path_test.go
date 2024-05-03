//spellchecker:words testutil
package testutil

//spellchecker:words testing
import (
	"os"
	"testing"
)

func TestToOSPath(t *testing.T) {
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
			if got := ToOSPath(tt.name); got != tt.want {
				t.Errorf("ToOSPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
