package utils

// SliceContainsAny checks if a slice contains any of a certain set of needles
func SliceContainsAny(haystack []string, needles ...string) bool {
	for _, a := range haystack {
		for _, needle := range needles {
			if a == needle {
				return true
			}
		}
	}
	return false
}
