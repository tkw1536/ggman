package path

import (
	"testing"
)

func TestToOSPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"hello/world", "hello" + pathSeperator + "world"},
		{"", ""},
		{"./", "." + pathSeperator},
		{"hello/../world", "hello" + pathSeperator + ".." + pathSeperator + "world"},
		{"/root", defaultVolumePrefix + pathSeperator + "root"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToOSPath(tt.name); got != tt.want {
				t.Errorf("ToOSPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
