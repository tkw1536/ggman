package path

import (
	"testing"

	"github.com/tkw1536/ggman/internal/testutil"
)

func TestGoesUp(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		// relative
		{"sub", false},
		{"sub/", false},
		{"../other", true},
		{"", false},
		{"..", true},
		{"./b/c", false},
		{"./../.", true},

		//absolute
		{"/sub", false},
		{"/sub/", false},
		{"/../other", false},
		{"/", false},
		{"/..", false},
		{"/./b/c", false},
		{"/./../.", false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := GoesUp(testutil.ToOSPath(tt.path)); got != tt.want {
				t.Errorf("PathGoesUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			if got := Contains(testutil.ToOSPath(tt.parent), testutil.ToOSPath(tt.child)); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
