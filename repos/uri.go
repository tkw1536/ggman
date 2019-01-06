package repos

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/tkw1536/ggman/utils"
)

// We support three types of urls:

// 1. The long form:
// [scheme://][user[:password]@]hostname[:port]/path
// e.g. https://git:git@mydomain:1234/repo.git
// e.g. https://git@mydomain/example

// 2. The short form:
// [scheme://][user[:password]@]hostname:path
// e.g. mydomain:hello/world.git

// RepoURI represents the URI to a single repository
type RepoURI struct {
	Scheme   string
	User     string
	Password string
	HostName string
	Port     int
	Path     string
}

// NewRepoURI parses a new repo uri from a string
func NewRepoURI(s string) (repo *RepoURI, err error) {
	// trim off the scheme and make sure that it validates
	scheme, rest := utils.SplitBefore(s, "://")
	if scheme != "" && !validateScheme(scheme) {
		rest = scheme + "://" + rest
	}

	// trim off authentication
	auth, rest := utils.SplitBefore(rest, "@")
	user, password := utils.SplitAfter(auth, ":")

	// if we have a ':', we need to determine if it is a port or a path
	var hostname, path, sport string
	var port int
	if strings.ContainsRune(rest, ':') {
		// trim off the hostname
		hostname, path = utils.SplitBefore(rest, ":")

		// if we have a scheme, then we have to parse everything after ':' as a port
		if scheme != "" {
			sport, path = utils.SplitAfter(path, "/")
			if port, err = parsePort(sport); err != nil {
				path = sport + "/" + path
				port = 0
				err = nil
			}
		}

	} else {
		// the first part of the url is a hostname
		hostname, path = utils.SplitAfter(rest, "/")
	}

	repo = &RepoURI{scheme, user, password, hostname, port, path}

	return
}

// IsCanonicalURI checks if a given URI is in canonical form
// using a specific canonical specification
func IsCanonicalURI(s string, cspec string) bool {
	uri, err := NewRepoURI(s)
	if err != nil {
		return false
	}
	return uri.Canonical(cspec) == s
}

var schemeRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9+\\-\\.]*$")

// validateScheme validates a scheme
func validateScheme(scheme string) bool {
	return schemeRegex.MatchString(scheme)
}

// parsePort parses a valid port
func parsePort(portstring string) (port int, err error) {
	port, err = strconv.Atoi(portstring)
	if err == nil && port < 0 || port > 65535 {
		err = errors.New("Port number outside of valid range 0 <= port <= 65535")
	}
	return
}

var reSlash = regexp.MustCompile("/+")

// Components gets the Components of a Repo URI
func (rURI *RepoURI) Components() (parts []string) {

	// normalize the path
	path := (*rURI).Path
	path = utils.TrimPrefixWhile(utils.TrimSuffixWhile(path, "/"), "/")
	path = strings.TrimSuffix(path, ".git")
	path = utils.TrimSuffixWhile(path, "/")
	path = reSlash.ReplaceAllString(path, "/")

	// get the host and the username
	host := (*rURI).HostName
	user := (*rURI).User

	// split the path into parts
	if path != "" {
		parts = strings.Split(path, "/")
	} else {
		parts = []string{}
	}

	// prepend (host, [user]) (with user iff a valid user exists)
	if user != "" && user != "git" {
		parts = append([]string{host, user}, parts...)
	} else {
		parts = append([]string{host}, parts...)
	}

	return
}

var specReplace = regexp.MustCompile("[\\^\\%]")

// Canonical returns the canonical version of this URI given a canonical specification
// the canonical specification can contain any character, except for three special ones
// ^ -- replaced by the first un-used component of the URI
// % -- replaced by the second un-used component of the URI (commonly username)
// $ -- replaced by all remaining components in the URI joined with a '/'. Also stops all processing afterwards
// If $ does not exist in the cspec, it is assumed to be at the end of the cspec.
func (rURI *RepoURI) Canonical(cspec string) (canonical string) {
	// get the components of the URI
	components := rURI.Components()

	// split into prefix and suffix
	prefix, suffix := utils.SplitAfter(cspec, "$")

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
