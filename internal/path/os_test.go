package path

import "testing"

func TestToOSPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"hello/world", "hello" + Separator + "world"},
		{"", ""},
		{"./", "." + Separator},
		{"hello/../world", "hello" + Separator + ".." + Separator + "world"},
		{"/root", defaultVolumePrefix + Separator + "root"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToOSPath(tt.name); got != tt.want {
				t.Errorf("ToOSPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
