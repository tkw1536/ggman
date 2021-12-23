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
