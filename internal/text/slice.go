package text

// SliceContainsAny returns true iff at least one needle is contained in haystack.
func SliceContainsAny(haystack []string, needles ...string) bool {
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
func SliceMatchesAny(haystack []string, predicates ...func(string) bool) bool {
	for _, hay := range haystack {
		for _, predicate := range predicates {
			if predicate(hay) {
				return true
			}
		}
	}
	return false
}

// SliceEquals checks if the first and second slice are equal.
//
// Two slices are considered equal when all elements are equal.
// Unlike reflect.DeepEqual, this considers any slices of zero length equal.
func SliceEquals(first, second []string) bool {
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

// SliceCopy returns a copy of the provided slice.
func SliceCopy(slice []string) []string {
	if len(slice) == 0 {
		return nil
	}
	clone := make([]string, len(slice))
	copy(clone, slice)
	return clone
}
