package slice

import (
	"reflect"
	"testing"
)

func TestContainsAny(t *testing.T) {
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
			if got := ContainsAny(tt.args.haystack, tt.args.needles...); got != tt.want {
				t.Errorf("ContainsAny() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMatchesAny(t *testing.T) {

	// p makes one predicate for each element in haystack
	p := func(haystack ...string) (fs []func(string) bool) {
		for _, hay := range haystack {
			h := hay
			fs = append(fs, func(s string) bool {
				return s == h
			})
		}
		return
	}

	type args struct {
		haystack   []string
		predicates []func(string) bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"single needle contained in haystack", args{
			haystack:   []string{"a", "b", "c"},
			predicates: p("a"),
		}, true},
		{"single needle not contained in haystack", args{
			haystack:   []string{"a", "b", "c"},
			predicates: p("d"),
		}, false},
		{"haystack contains a single needle", args{
			haystack:   []string{"a", "b", "c"},
			predicates: p("f", "a", "e"),
		}, true},
		{"haystack contains no needle", args{
			haystack:   []string{"a", "b", "c"},
			predicates: p("d", "e", "f"),
		}, false},
		{"empty haystack", args{
			haystack:   nil,
			predicates: p("d", "e", "f"),
		}, false},
		{"empty needles", args{
			haystack:   []string{"a", "b", "c"},
			predicates: nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchesAny(tt.args.haystack, tt.args.predicates...); got != tt.want {
				t.Errorf("MatchesAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquals(t *testing.T) {
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
			if got := Equals(tt.args.first, tt.args.second); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	type args struct {
		slice []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"empty slice", args{nil}, nil},
		{"non-empty slice", args{[]string{"a", "b", "c"}}, []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Copy(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveZeros(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"remove from the nil slice", args{nil}, nil},
		{"remove from the empty array", args{[]string{}}, []string{}},
		{"remove from some places", args{[]string{"", "x", "y", "", "z"}}, []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveZeros(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveZeros() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkRemoveZeros(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RemoveZeros[struct{}](nil)
		RemoveZeros([]string{})
		RemoveZeros([]string{"", "x", "y", "", "z"})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"nil slice", args{nil}, nil},
		{"no duplicates", args{[]string{"a", "b", "c", "d"}}, []string{"a", "b", "c", "d"}},
		{"some duplicates", args{[]string{"b", "c", "c", "d", "a", "b", "c", "d"}}, []string{"a", "b", "c", "d"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicates(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}
