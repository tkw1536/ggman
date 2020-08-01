package util

import "testing"

func TestSliceContainsAny(t *testing.T) {
	type args struct {
		haystack []string
		needles  []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"single needle contained in haystack", args{
			haystack: []string{"a", "b", "c"},
			needles:  []string{"a"},
		}, true},
		{"single needle not contained in haystack", args{
			haystack: []string{"a", "b", "c"},
			needles:  []string{"d"},
		}, false},
		{"haystack contains a single needle", args{
			haystack: []string{"a", "b", "c"},
			needles:  []string{"f", "a", "e"},
		}, true},
		{"haystack contains no needle", args{
			haystack: []string{"a", "b", "c"},
			needles:  []string{"d", "e", "f"},
		}, false},
		{"empty haystack", args{
			haystack: nil,
			needles:  []string{"d", "e", "f"},
		}, false},
		{"empty needles", args{
			haystack: []string{"a", "b", "c"},
			needles:  nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceContainsAny(tt.args.haystack, tt.args.needles...); got != tt.want {
				t.Errorf("SliceContainsAny() = %v, want %v", got, tt.want)
			}
		})
	}
}
