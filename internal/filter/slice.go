// Package filter provides slice filtering functions
package filter

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Inplace filters a slice in-place using pred.
//
// This means that it creates a new slice using the same backing array of slice that only contains values
// for which pred(v) returns True.
//
// This trivially invalidates the old value of slice.
// As such the result value should typically be assigned to the input value.
// For example:
//
//	s = Inplace(s, pred)
func Inplace[T any](slice []T, pred func(T) bool) []T {
	// check that we have a predicate!
	if pred == nil {
		panic("FilterI: pred is nil")
	}

	// when the slice is nil, we have nothing to filter
	if slice == nil {
		return nil
	}

	// create a new result slice
	result := slice[:0]
	for _, v := range slice {
		if !pred(v) {
			continue
		}
		result = append(result, v)
	}

	// if we still have some leftover elements we need to prevent memory leaks
	// so zero out the rest of the slice.
	if len(result) < len(slice) {
		// outer if is an optimization to prevent allocation when not needed!
		var zeroT T
		for i := len(result); i < len(slice); i++ {
			slice[i] = zeroT
		}
	}

	// and return the result slice!
	return result
}

// RemoveZeros removes zero values from s in-place.
//
// This trivially invalidates the old value of slice.
// As such the result value should typically be assigned to the input value.
// For example:
//
//	s = RemoveZeros(s)
func RemoveZeros[T comparable](s []T) []T {
	var zeroT T
	return Inplace(s, func(v T) bool { return v != zeroT })
}

// RemoveDuplicates removes duplicates in s.
// As a side effect, elements in s are also ordered.
//
// This function will invalidate the previous value of s.
//
// It is recommended to store the return value of this function in the original variable.
// The call should look something like:
//
//	s = RemoveDuplicates(s)
func RemoveDuplicates[T constraints.Ordered](s []T) []T {
	slices.Sort(s)
	return slices.Compact(s)
}
