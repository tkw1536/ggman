package path

import (
	"testing"

	"github.com/tkw1536/ggman/internal/testutil"
)

func TestContains(t *testing.T) {
	tests := []struct {
		parent string
		child  string
		want   bool
	}{
		{"/root/", "/root/child/", true},
		{"/root/", "/other/", false},
		{"/root/", "/root/", true},

		{"./root/", "./root/child/", true},
		{"./root/", "./other/", false},
		{"./root/", "./root/", true},
	}
	for _, tt := range tests {
		t.Run(tt.parent+" contains "+tt.child, func(t *testing.T) {
			if got := HasChild(testutil.ToOSPath(tt.parent), testutil.ToOSPath(tt.child)); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
