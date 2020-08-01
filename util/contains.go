// Package util contains various utility functions used by ggman
package util

// SliceContainsAny returns true if at least one needle is contained in haystack.
// Otherwise returns false.
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
