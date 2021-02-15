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

func TestSliceEquals(t *testing.T) {
	type args struct {
		first  []string
		second []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{

		{"equality of nil slices", args{nil, nil}, true},
		{"equality of nil and empty slice", args{nil, []string{}}, true},
		{"equality of empty and nil slice", args{[]string{}, nil}, true},
		{"equality of empty slices", args{[]string{}, []string{}}, true},

		{"inequality of empty and full slice", args{nil, []string{"a"}}, false},
		{"inequality of full and empty slice", args{[]string{"a"}, nil}, false},

		{"equality of full slices", args{[]string{"a"}, []string{"a"}}, true},
		{"equality of full slices (2)", args{[]string{"a", "b", "c"}, []string{"a", "b", "c"}}, true},

		{"inequality of full slices", args{[]string{"a"}, []string{"a", "b", "c"}}, false},
		{"inequality of full slices (2)", args{[]string{"a", "b", "c"}, []string{"a"}}, false},
		{"inequality of full slices (3)", args{[]string{"a", "b", "c"}, []string{"a", "d", "c"}}, false},
		{"inequality of full slices (4)", args{[]string{"a", "d", "c"}, []string{"a", "b", "c"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceEquals(tt.args.first, tt.args.second); got != tt.want {
				t.Errorf("SliceEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}
