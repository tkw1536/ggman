package env

//spellchecker:words strings github ggman internal split pkglib collection text jessevdk flags
import (
	"strings"

	"github.com/tkw1536/ggman/internal/split"
	"github.com/tkw1536/ggman/internal/url"
	"github.com/tkw1536/pkglib/collection"
	"github.com/tkw1536/pkglib/text"

	"github.com/jessevdk/go-flags"
)

//spellchecker:words mydomain

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

var (
	_ flags.Unmarshaler = (*URL)(nil)
)

// UnmarshalFlag implements the flags.Unmarshaler interface.
func (u *URL) UnmarshalFlag(value string) error {
	*u = ParseURL(value)
	return nil
}

var windowsReplacer = strings.NewReplacer("\\", "/")

// ParseURL parses a string into a repo URL.
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
// ParseURL always succeeds.
// This can lead to unexpected parses of URLs when e.g. a port is specified incorrectly.
//
// For windows compatibility, '\\' is replaced by '/' in the input string.
func ParseURL(s string) (repo URL) {
	// normalize for windows
	s = windowsReplacer.Replace(s)

	// split off the url scheme
	repo.Scheme, s = url.SplitURLScheme(s)

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
	// First split the path into components split by '/'
	// Adding the two '//' makes sure that there is enough space in the slice.
	components := collection.KeepFunc(strings.Split(url.Path, "/"), func(s string) bool { return s != "" })

	// remove the last component that has a '.git' in it
	if last := len(components) - 1; last >= 0 {
		components[last] = strings.TrimSuffix(components[last], ".git")

		// if we had a '/' before the .git, remove it
		if components[last] == "" {
			components = components[:last]
		}
	}

	// Prepend the hostname and user to the components.
	if url.User == "" || url.User == "git" {
		components = append(components, "") // add space to the hostname
		copy(components[1:], components)
		components[0] = url.HostName
	} else {
		components = append(components, "", "") // add two spaces to the front of the array
		copy(components[2:], components)
		components[0] = url.HostName
		components[1] = url.User
	}

	return components
}

// Canonical returns the canonical version of this URI given a canonical specification
// the canonical specification can contain any character, except for three special ones
// ^ -- replaced by the first un-used component of the URI
// % -- replaced by the second un-used component of the URI (commonly username)
// $ -- replaced by all remaining components in the URI joined with a '/'. Also stops all processing afterwards.
// If $ does not exist in the cSpec, it is assumed to be at the end of the cSpec.
func (url URL) Canonical(cSpec string) (canonical string) {
	var builder strings.Builder

	components := url.Components()                // get the components of the URI
	prefix, suffix := split.AfterRune(cSpec, '$') // split into mod-able and static part

	for i, r := range prefix {
		// no more components left
		// => we can immediately exit the loop
		if len(components) == 0 {
			builder.WriteString(prefix[i:])
			break
		}

		switch r {
		case '%':
			// insufficient components.
			if len(components) < 2 {
				builder.WriteRune(r)
				break /* switch */
			}

			// write the second component
			builder.WriteString(components[1])
			components[1] = components[0]
			components = components[1:]
		case '^':
			// write the first component
			builder.WriteString(components[0])
			components = components[1:]
		default:
			builder.WriteRune(r)
		}
	}

	// add all the components to replace the '$'
	_, _ = text.Join(&builder, components, "/") // ignore cause this should never fail

	// add the suffix
	builder.WriteString(suffix)

	return builder.String()
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
