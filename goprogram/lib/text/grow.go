package text

import "io"

// Grow calls w.Grow() when w provides a Grow() method
func Grow(w io.Writer, n int) {
	if grower, canGrow := w.(interface {
		Grow(length int)
	}); canGrow {
		grower.Grow(n)
	}
}
