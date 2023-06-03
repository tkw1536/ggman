package env

import (
	"regexp"
	"strings"

	"github.com/tkw1536/ggman/internal/split"
	"github.com/tkw1536/ggman/internal/url"
	"github.com/tkw1536/pkglib/collection"
)

// URL represents a URL to a single git repository.
//
// A URL consists of different parts, and is typically created by using ParseURL.
//
// URLs can be both concrete (that is representing a single repository) or a "pattern" matching multiple URLS.
// The implementation does not distinguish between either case.
// For pattern matching see the Match function.
type URL struct {
	Scheme string // e.g. "ssh"

	User     string // e.g. "git"
	Password string // e.g. "p@ssw0rd"

	HostName string // e.g. "github.com"
	Port     uint16 // e.g. 2222

	Path string // e.g. "hello/world.git"
}

func (u *URL) UnmarshalFlag(value string) error {
	*u = ParseURL(value)
	return nil
}

// ParseURL parses a URL without any namespaces.
// See ParseURLNamespace.
func ParseURL(s string) (repo URL) {
	s = normPath(s)

	// Trim off a leading scheme (as separated by '://') and (if it is valid) store it.
	scheme, rest := split.Before(s, "://")
	if url.IsValidURLScheme(scheme) {
		repo.Scheme = scheme
		s = rest
	}

	// Next, we split of the authentication if we have an '@' sign.
	// Technically the if { } clause isn't required, the code will work fine without it.
	// However most URLs will not have an '@' sign, so we can save allocating an extra variable and the function call.
	if at := strings.IndexRune(s, '@'); at > 0 {
		var auth string
		auth, s = s[:at], s[at+1:]
		repo.User, repo.Password = split.AfterRune(auth, ':')
	}

	// Finally, we check if the remainder contains a ':'.
	// If it does, we have to figure out if it is of the form hostname:port or hostname:path.
	// The second form is only allowed if we have some kind of scheme.
	// If there is no ':', we can straightforwardly split after the first '/'
	if colon := strings.IndexRune(s, ':'); colon >= 0 {
		repo.HostName, s = s[:colon], s[colon+1:]

		// if we have a scheme, then we have to parse everything after ':' as a port.
		// This only works if the port is valid.
		if repo.Scheme != "" {
			var err error
			port, rest := split.AfterRune(s, '/')
			if repo.Port, err = url.ParsePort(port); err == nil {
				s = rest
			}
		}

		repo.Path = s
		return
	}

	repo.HostName, repo.Path = split.AfterRune(s, '/')
	return
}

// ParseURLContext parses a string into a repo URL, while expanding the provided namespaces.
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
//
// A namespace may be prepended to the URL using the form 'namespace:[//]'
// When a namespace is encountered, the remainder of the URL is appended to the expansion of the namespace.
// Then the URL parsing is restarted.
//
// Namespaces may not be nested, meaning only one expansion of namespaces is done.
// This is to prevent infinite loops during parsing.
//
// ParseURLContext always succeeds.
// This can lead to unexpected parses of URLs when e.g. a port is specified incorrectly.
//
// For windows compatibility, '\\' is replaced by '/' in the input string.
func ParseURLContext(s string, namespaces map[string]string) (repo URL) {
	s = normPath(s)

	// substitute the namespaces beforehand
	// check the long form before the short form
	for ns, expansion := range namespaces {
		if v, ok := strings.CutPrefix(s, ns+"://"); ok {
			s = expansion + v
			break
		}

		if v, ok := strings.CutPrefix(s, ns+":"); ok {
			s = expansion + v
			break
		}
	}

	return ParseURL(s)
}

var windowsReplacer = strings.NewReplacer("\\", "/")

// normPath normalizes a path to use the same separator on windows systems
func normPath(s string) string {
	return windowsReplacer.Replace(s)
}

// IsLocal checks if this URL looks like a local URL.
// A URL is considered local if it uses the "file" scheme, or the scheme is empty and the hostname is one of ".", ".." or "".
func (url URL) IsLocal() bool {
	return url.Scheme == "file" || (url.Scheme == "" && (url.HostName == "." || url.HostName == ".." || url.HostName == ""))
}

// Components gets the components of a URL
//
// Components of the URL are the hostname, the username and components of the path.
// Empty components are ignored.
// Furthermore a username 'git' as well as a trailing suffix of '.git' are ignored as well.
func (url URL) Components() []string {
	// First split the path into components split by '/'.
	// and remove a '.git' from the last part.
	parts := collection.Filter(strings.Split(url.Path, "/"), func(s string) bool { return s != "" })
	lastPart := len(parts) - 1
	if lastPart >= 0 {
		parts[lastPart] = strings.TrimSuffix(parts[lastPart], ".git")

		// if we had a '/' before the .git, remove it
		if parts[lastPart] == "" {
			parts = parts[:lastPart]
		}
	}

	// Now prepend the hostname and user (unless it is 'git' or missing)
	components := make([]string, 1, 2+len(parts))
	components[0] = url.HostName
	if url.User != "" && url.User != "git" {
		components = append(components, url.User)
	}
	return append(components, parts...)
}

var specReplace = regexp.MustCompile(`[\^\%]`)

// Canonical returns the canonical version of this URI given a canonical specification
// the canonical specification can contain any character, except for three special ones
// ^ -- replaced by the first un-used component of the URI
// % -- replaced by the second un-used component of the URI (commonly username)
// $ -- replaced by all remaining components in the URI joined with a '/'. Also stops all processing afterwards
// If $ does not exist in the cspec, it is assumed to be at the end of the cspec.
func (url URL) Canonical(cspec string) (canonical string) {
	// get the components of the URI
	components := url.Components()

	// split into prefix and suffix
	prefix, suffix := split.AfterRune(cspec, '$')

	prefix = specReplace.ReplaceAllStringFunc(prefix, func(s string) string {
		// if everything is empty, return the string as is
		if len(components) == 0 {
			return s
		}

		// replace by the first component
		if s == "^" {
			defer func() { components = components[1:] }()
			return components[0]
		}

		// else we want to replace by the second component
		// so we need to make sure we have that many
		if len(components) < 2 {
			return s
		}

		// do the replacement
		defer func() { components = append(components[:1], components[2:]...) }()
		return components[1]
	})

	// add the remaining components
	return prefix + strings.Join(components, "/") + suffix
}

// CanonicalWith returns the canonical url given a set of lines
func (url URL) CanonicalWith(lines CanFile) (canonical string) {
	var pat PatternFilter
	for _, line := range lines {
		pat.Set(line.Pattern)
		if pat.MatchesURL(url) {
			return url.Canonical(line.Canonical)
		}
	}

	return
}

// ComponentsOf returns the components of the URL in s.
// It is a convenience wrapper for ParseURL(s).Components().
func ComponentsOf(s string) []string {
	return ParseURL(s).Components()
}
