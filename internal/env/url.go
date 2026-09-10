package env

//spellchecker:words strconv strings ggman internal parseurl split pkglib text
import (
	"strconv"
	"strings"

	"go.tkw01536.de/ggman/internal/parseurl"
	"go.tkw01536.de/ggman/internal/split"
	"go.tkw01536.de/pkglib/text"
)

//spellchecker:words mydomain recvcheck

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

	HostName string // e.g. "gitforge.example"
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

// ParseURLAndForgeReference is like [ParseURL] followed by [URL.SplitForgeReference].
//
// The returned URL contains the remaining path after splitting off the tree reference.
func ParseURLAndForgeReference(s string) (url URL, suffix, ref, relative string) {
	url = ParseURL(s)
	url.Path, suffix, ref, relative = url.SplitForgeReference()
	return
}

// IsLocal checks if this URL looks like a local URL.
// A URL is considered local if it uses the "file" scheme, or the scheme is empty and the hostname is one of ".", ".." or "".
func (url URL) IsLocal() bool {
	return url.Scheme == "file" || (url.Scheme == "" && (url.HostName == "." || url.HostName == ".." || url.HostName == ""))
}

// IsForgeURL checks if this URL looks like a URL copied from a web browser.
// A URL is considered a web URL if it has a scheme of "https" or "http", or the scheme is empty and the hostname is not one of ".", ".." or "".
func (url URL) IsWebURL() bool {
	switch strings.ToLower(url.Scheme) {
	case "https", "http":
		return true
	case "":
		return url.HostName != "" && url.HostName != "." && url.HostName != ".."
	default:
		return false
	}
}

// SplitForgeReference splits a forge-like tree reference into a clean path, a suffix, a source control reference and a relative path.
//
// For example the URL "https://gitforge.example/user/repo/tree/ref/path?tab=README" would be split into:
//
//	("user/repo", "?tab=README", "ref", "path")
//
// The reference consists of:
// - a cleaned up path ("path") without the tree reference
// - a suffix ("suffix") containing any stray query or fragment identifiers
// - a source control reference ("ref")
// - a relative path from the root of the repository ("relative")
//
// Parsing is domain-agnostic and based only on path segment shape (no hostname checks).
// Path segments are matched exactly, so names like "restructure" or "treehouse" are not treated as markers.
//
// Supported structural markers (first match wins, left to right):
//
//   - a "tree" or "blob" segment followed by a non-empty ref segment
//     (GitHub / Gitea / Forgejo / Sourcehut-style, e.g. "~user/repo/tree/main/src")
//   - a "src" segment followed by "branch" or "tag", then a non-empty ref segment
//     (Gitea / Forgejo, e.g. "user/repo/src/branch/main/src" or "user/repo/src/tag/v1.0")
//
// The ref is always exactly one path segment after the marker; remaining segments form the relative path.
// Multi-segment branch names must appear as a single segment (e.g. "feature%2Fx").
//
// Additional heuristics:
//
// - Non-web URLs, as reported by [URL.IsWebURL], are not considered forge URLs, and return their original path unmodified, along with empty suffixes, ref and relative.
// - Trailing "_" or "-" segments immediately before the marker are removed (e.g. Gitea "/_/" and GitLab "/-/").
// - A trailing query or fragment identifier of the URL is split off and returned in the suffix instead.
//
// If no tree reference is found, path is the original path (without query or fragment), and ref and
// relative are empty.
func (url URL) SplitForgeReference() (path, suffixes, ref, relative string) {
	if !url.IsWebURL() {
		return url.Path, "", "", ""
	}

	// Split off a trailing query or fragment identifier.
	if i := strings.IndexAny(url.Path, "?#"); i >= 0 {
		url.Path, suffixes = url.Path[:i], url.Path[i:]
	}

	components := parseurl.SplitNonEmpty(url.Path, '/', nil)
	cut, refStart := -1, -1
	for i := range components {
		switch {
		case (components[i] == "tree" || components[i] == "blob") && i+1 < len(components):
			cut, refStart = i, i+1
		case components[i] == "src" && i+2 < len(components) && (components[i+1] == "branch" || components[i+1] == "tag"):
			cut, refStart = i, i+2
		default:
			continue
		}
		break
	}
	if cut < 0 {
		return url.Path, suffixes, "", ""
	}

	prefix := components[:cut]
	for len(prefix) > 0 {
		last := prefix[len(prefix)-1]
		if last != "_" && last != "-" {
			break
		}
		prefix = prefix[:len(prefix)-1]
	}

	ref = components[refStart]
	relative = strings.Join(components[refStart+1:], "/")
	return strings.Join(prefix, "/"), suffixes, ref, relative
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
// ! -- consumes a component without using it
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
		case '!':
			// drop the component of the URL
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
