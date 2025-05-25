//spellchecker:words parseurl
package parseurl

// SplitNonEmpty is like [strings.Split] except that empty parts are skipped.
// Furthermore, instead of directly returning the result it is appended to buf.
// The appended buf is returned.
func SplitNonEmpty(s string, sep rune, buf []string) []string {
	if buf == nil {
		count := CountNonEmptySplit(s, sep)
		buf = make([]string, 0, count)
	}

	var (
		lastWasSep    = true // was the last rune a separator?
		lastPartStart int    // where did the last part start?
	)

	for i, c := range s {
		if lastWasSep {
			lastPartStart = i
		}

		isSep := c == sep

		if !lastWasSep && isSep {
			buf = append(buf, s[lastPartStart:i])
		}

		lastWasSep = isSep
	}

	if !lastWasSep {
		buf = append(buf, s[lastPartStart:])
	}

	return buf
}

// CountNonEmptySplit returns the number of non-empty
// non-overlapping contiguous sequences of runes of s
// that do not include sep.
//
// This corresponds to the length of [SplitNonEmpty]
// with a nil buffer.
func CountNonEmptySplit(s string, sep rune) (count int) {
	lastWasSep := true // was the last rune a separator?

	for _, c := range s {
		isSep := c == sep

		if !isSep && lastWasSep {
			count++
		}

		lastWasSep = isSep
	}
	return count
}
