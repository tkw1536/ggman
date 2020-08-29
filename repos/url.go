package repos

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/tkw1536/ggman/util"
)

// URL represents the URL to a single git repository.
//
// URLs can be both concrete and a pattern matching a list of URLs.
// The implementation does not distinguish between either case.
//
// For supported syntaxes of URLs, see ParseRepoURL.
type URL struct {
	Scheme   string
	User     string
	Password string
	HostName string
	Port     uint16
	Path     string
}

// ParseRepoURL parses a new repo uri from a string.
//
// We support two types of urls:
//
// 1. The long form:
// [scheme://][user[:password]@]hostname[:port]/path
// e.g. https://git:git@mydomain:1234/repo.git
// e.g. https://git@mydomain/example
//
// 2. The short form:
// [scheme://][user[:password]@]hostname:path
// e.g. mydomain:hello/world.git
func ParseRepoURL(s string) (repo URL) {
	// we sometimes have to restore 's' to what it was before
	// we could do this using string concatination, but that is slow.
	// so we store it in a temporary variable called 'oldS'

	// If the URL starts with a scheme, we try to trim it off here.
	// If we find out that this is not a valid scheme,
	// We instead revert s to what it was before we split it off.
	oldS := s
	repo.Scheme, s = util.SplitBefore(s, "://")
	if repo.Scheme != "" && !IsValidURLScheme(repo.Scheme) {
		s = oldS
		repo.Scheme = ""
	}

	// Next, we split of the authentication if we have an '@' sign.
	// Technically the if { } clause isn't required, the code will work even without it.
	// However most URLs will not have an '@' sign, so we can save allocating an extra variable and the function call.
	if strings.ContainsRune(s, '@') {
		var auth string
		auth, s = util.SplitBefore(s, "@")
		repo.User, repo.Password = util.SplitAfter(auth, ":")
	}

	// Finally, we cherck if the remainder contains a ':'.
	// If it does, we have to figure out if it is of the form hostname:port or hostname:path.
	// The second form is only allowed if we have some kind of scheme.
	// If there is no ':', we can straightforwardly split after the first '/'
	if strings.ContainsRune(s, ':') {
		repo.HostName, s = util.SplitBefore(s, ":")

		// if we have a scheme, then we have to parse everything after ':' as a port
		// However we need to make sure that port is
		if repo.Scheme != "" {
			var portString string
			oldS := s
			portString, s = util.SplitAfter(s, "/")
			var err error
			if repo.Port, err = ParsePort(portString); err != nil {
				s = oldS
			}
		}

		repo.Path = s
	} else {
		repo.HostName, repo.Path = util.SplitAfter(s, "/")
	}

	return
}

// IsValidURLScheme checks if the string s reprents a valid URL scheme.
//
// A valid url scheme is one that matches the regex [a-zA-Z][a-zA-Z0-9+\-\.]*.
func IsValidURLScheme(s string) bool {
	// An obvious implementation of this function would be:
	//   regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9+\-\.]*$").MatchString(s)
	// However such an implementation would be relatively slow when called a lot of times.
	// Instead we directly build code that implements this RegEx.
	//
	// For this we first check that the string is non-empty.
	// Then we check that each character is either alphanumeric or a '+', '-' or '.'
	// Finally we check that the first character is not a '+' or '-'.

	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if !(('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')) && // more efficient version of !@unicode.isLetter(r)
			r != '+' &&
			r != '-' &&
			r != '.' {
			return false
		}

	}
	return s[0] != '+' && s[0] != '-' && s[0] != '.'
}

var errNoPlus = errors.New("ParsePort: s may not start with '+'")
var errInvalidRange = errors.New("ParsePort: s out of range")

// ParsePort parses the string s into a valid port.
// When a port can not be parsed, returns 0 and an error.
func ParsePort(s string) (uint16, error) {
	// This function could use strconv.ParseUint(s, 10, 16) and then do a conversion on the result.
	// Instead we call strconv.AtoI and check that the result is in the right range.
	//
	// In Benchmarking it turned out that this implementation is about 33% faster for common port usecases.
	// This is likely because strconv.AtoI has a 'fast' path for small integers, and most port candidates are small enough to fit.

	if len(s) > 0 && s[0] == '+' {
		return 0, errNoPlus
	}
	port, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if port < 0 || port > 65535 {
		return 0, errInvalidRange
	}
	return uint16(port), nil
}
