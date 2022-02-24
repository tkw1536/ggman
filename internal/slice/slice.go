// Package slice contains utility function for slices
package slice

// ContainsAny returns true iff at least one needle is contained in haystack.
func ContainsAny[T comparable](haystack []T, needles ...T) bool {
	for _, hay := range haystack {
		for _, needle := range needles {
			if hay == needle {
				return true
			}
		}
	}
	return false
}

// SliceContainsAny returns true iff at least one predicate matches any of haystack.
func MatchesAny[T comparable](haystack []T, predicates ...func(T) bool) bool {
	for _, hay := range haystack {
		for _, predicate := range predicates {
			if predicate(hay) {
				return true
			}
		}
	}
	return false
}

// Equals checks if the first and second slice are equal.
//
// Two slices are considered equal when all elements are equal and occur in the same order.
// Unlike reflect.DeepEqual, this considers any slices of zero length equal.
func Equals[T comparable](first, second []T) bool {
	if len(first) != len(second) {
		return false
	}
	for idx, value := range first {
		if value != second[idx] {
			return false
		}
	}
	return true
}

// Copy returns a copy of the provided slice.
func Copy[T any](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}
	clone := make([]T, len(slice))
	copy(clone, slice)
	return clone
}

// FilterI filters a slice in-place using pred.
//
// This means that it creates a new slice using the same backing array of slice that only contains values
// for which pred(v) returns True.
//
// This trivially invalidates the old value of slice.
// As such the result value should typically be assigned to the input value.
// For example:
//
//   s = FilterI(s, pred)
func FilterI[T any](slice []T, pred func(T) bool) []T {
	// create a new result slice
	result := slice[:0]
	for _, v := range slice {
		if !pred(v) {
			continue
		}
		result = append(result, v)
	}

	// zero out all the now unused values of T
	// this prevents memory leaks
	var zeroT T
	for i := len(result); i < len(slice); i++ {
		slice[i] = zeroT
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
//   s = RemoveZeros(s)
func RemoveZeros[T comparable](s []T) []T {
	var zeroT T
	return FilterI(s, func(v T) bool { return v != zeroT })
}

// RemoveDuplicates removes duplicates in s.
// As a side effect, elements in s are also ordered.
//
// This function will invalidate the previous value of s.
//
// It is recommended to store the return value of this function in the original variable.
// The call should look something like:
//
//  s = RemoveDuplicates(s)
//
func RemoveDuplicates[T Ordered](s []T) []T {
	if len(s) == 0 {
		return s
	}

	// adapted from https://github.com/golang/go/wiki/SliceTricks#in-place-deduplicate-comparable
	Sort(s)

	j := 0
	for i := 1; i < len(s); i++ {
		if s[j] == s[i] {
			continue
		}
		j++

		s[j] = s[i]
	}

	return s[:j+1]
}
