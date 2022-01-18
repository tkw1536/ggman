package path

import (
	"testing"
)

func TestGoesUp(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"sub/", false},
		{"../other", true},
		{"", false},
		{"..", true},
		{"./a/b/c", false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := GoesUp(tt.path); got != tt.want {
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
			if got := Contains(tt.parent, tt.child); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
