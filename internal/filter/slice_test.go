package filter

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInplace(t *testing.T) {
	type args struct {
		slice []int
		pred  func(int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "nil list",
			args: args{slice: nil, pred: func(i int) bool { panic("never reached") }},
			want: nil,
		},

		{
			name: "true filter",
			args: args{slice: []int{0, 1, 2, 3}, pred: func(i int) bool { return true }},
			want: []int{0, 1, 2, 3},
		},
		{
			name: "false filter",
			args: args{slice: []int{0, 1, 2, 3}, pred: func(i int) bool { return false }},
			want: []int{},
		},

		{
			name: "even filter",
			args: args{slice: []int{0, 1, 2, 3}, pred: func(i int) bool { return i%2 == 0 }},
			want: []int{0, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Inplace(tt.args.slice, tt.args.pred); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Inplace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleInplace() {

	// create a slice and filter it in-place!
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	filtered := Inplace(slice, func(i int) bool { return i%2 == 0 })

	// the filtered slice is as we would expect
	fmt.Println(filtered)

	// the original slice has been invalidated, elements not used have been zeroed out.
	fmt.Println(slice)

	// because we filtered in place, slice[0:6] refers to the same underlying array as filtered[0:6]
	// we show this by setting all of slice and printing it again
	slice[0] = -1
	slice[1] = -1
	slice[2] = -1
	slice[3] = -1
	slice[4] = -1
	slice[5] = -1
	fmt.Println(filtered)

	// Normally one would just
	//  slice = FilterI(slice, ...)
	// to prevent accidentally leaking memory.

	// Output: [0 2 4 6 8]
	// [0 2 4 6 8 0 0 0 0 0]
	// [-1 -1 -1 -1 -1]
}

func BenchmarkInplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Inplace([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, func(i int) bool { return i%2 == 0 })
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
