package slice

import "sort"

// Ordered is a type constraint that matches any ordered type.
// An ordered type is one that supports the <, <=, >, and >= operators.
type Ordered interface {
	// adapted from https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// Sort sorts slice ascendingly according to the < operator.
func Sort[T Ordered](s []T) {
	sort.Sort(slice[T](s))
}

type slice[T Ordered] []T

func (s slice[T]) Len() int           { return len(s) }
func (s slice[T]) Less(i, j int) bool { return s[i] < s[j] }
func (s slice[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// SortFunc sorts s ascendingly according to the less function.
func SortFunc[T any](s []T, less func(i, j T) bool) {
	sort.Sort(slicefunc[T]{
		data: s,
		less: less,
	})
}

type slicefunc[T any] struct {
	data []T
	less func(i, j T) bool
}

func (s slicefunc[T]) Len() int           { return len(s.data) }
func (s slicefunc[T]) Less(i, j int) bool { return s.less(s.data[i], s.data[j]) }
func (s slicefunc[T]) Swap(i, j int)      { s.data[i], s.data[j] = s.data[j], s.data[i] }
