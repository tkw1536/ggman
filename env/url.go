package env

//spellchecker:words strconv strings github ggman internal parseurl split pkglib text jessevdk flags
import (
	"strconv"
	"strings"

	"github.com/tkw1536/ggman/internal/parseurl"
	"github.com/tkw1536/ggman/internal/split"
	"go.tkw01536.de/pkglib/text"

	"github.com/jessevdk/go-flags"
)

//spellchecker:words mydomain nolint recvcheck

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

// String returns a string representation of this URL.
// It is a best guess of what was used to parse the URL, but not guaranteed to be so.
func (url URL) String() string {
	var builder strings.Builder

	if url.Scheme != "" {
		builder.WriteString(url.Scheme)
		builder.WriteString("://")
	}

	if url.User != "" {
		builder.WriteString(url.User)
		if url.Password != "" {
			builder.WriteString(":")
			builder.WriteString(url.Password)
		}
		builder.WriteString("@")
	}

	builder.WriteString(url.HostName)
	if url.Port != 0 {
		builder.WriteString(":")
		builder.WriteString(strconv.FormatUint(uint64(url.Port), 10))
	}

	if url.Path != "" {
		if url.Scheme != "" || url.IsLocal() {
			builder.WriteString("/")
		} else {
			builder.WriteString(":")
		}
		builder.WriteString(url.Path)
	}

	return builder.String()
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

// ParseURL parses a string into a URL.
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
func ParseURL(s string) (url URL) {
	// In this function we use url.Path as scratch space.
	// we keep splitting off parts as we parse them
	url.Path = windowsReplacer.Replace(s)                 // normalize for windows!
	url.Scheme, url.Path = parseurl.SplitScheme(url.Path) // split off the scheme

	// split off authentication, if any.
	if at := strings.IndexRune(url.Path, '@'); at >= 0 {
		url.User, url.Path = url.Path[:at], url.Path[at+1:]
		url.User, url.Password = split.AfterRune(url.User, ':')
	}

	colon := strings.IndexRune(url.Path, ':')
	if colon < 0 {
		url.HostName, url.Path = split.AfterRune(url.Path, '/')
		return
	}

	// we have the form "hostname:port/path" or "hostname:path".
	// the former is only valid if we have a scheme.
	url.HostName, url.Path = url.Path[:colon], url.Path[colon+1:]
	if url.Scheme == "" {
		return
	}

	// split off a valid port from the path.
	if slash := strings.IndexRune(url.Path, '/'); slash >= 0 {
		var err error
		url.Port, err = parseurl.ParsePort(url.Path[:slash])
		if err == nil {
			url.Path = url.Path[slash+1:]
		}
	}

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
	hasUser := url.User != "" && url.User != "git"

	count := parseurl.CountNonEmptySplit(url.Path, '/') + 1
	if hasUser {
		count += 1
	}

	components := make([]string, 1, count)
	components[0] = url.HostName
	if hasUser {
		components = append(components, url.User)
	}

	components = parseurl.SplitNonEmpty(url.Path, '/', components)

	// remove trailing '.git'
	if last := len(components) - 1; last >= 0 {
		components[last] = strings.TrimSuffix(components[last], ".git")

		if components[last] == "" {
			components = components[:last]
		}
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
	if cSpec == "$$" {
		return url.String()
	}
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
// If no pattern matches, return the best-guess original url.
func (url URL) CanonicalWith(lines CanFile) (canonical string) {
	var pat PatternFilter
	for _, line := range lines {
		pat.Set(line.Pattern)
		if pat.MatchesURL(url) {
			return url.Canonical(line.Canonical)
		}
	}

	return url.String()
}

// ComponentsOf returns the components of the URL in s.
// It is a convenience wrapper for ParseURL(s).Components().
func ComponentsOf(s string) []string {
	return ParseURL(s).Components()
}
