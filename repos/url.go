package repos

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/tkw1536/ggman/util"
)

// RepoURL represents the URI to a single repository
type RepoURL struct {
	Scheme   string
	User     string
	Password string
	HostName string
	Port     int
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
func ParseRepoURL(s string) (repo *RepoURL) {
	// trim off the scheme and make sure that it validates
	scheme, rest := util.SplitBefore(s, "://")
	if scheme != "" && !validateScheme(scheme) {
		rest = scheme + "://" + rest
	}

	// trim off authentication
	auth, rest := util.SplitBefore(rest, "@")
	user, password := util.SplitAfter(auth, ":")

	// if we have a ':', we need to determine if it is a port or a path
	var hostname, path, sport string
	var port int
	if strings.ContainsRune(rest, ':') {
		// trim off the hostname
		hostname, path = util.SplitBefore(rest, ":")

		// if we have a scheme, then we have to parse everything after ':' as a port
		if scheme != "" {
			sport, path = util.SplitAfter(path, "/")
			var err error
			if port, err = parsePort(sport); err != nil {
				path = sport + "/" + path
			}
		}

	} else {
		// the first part of the url is a hostname
		hostname, path = util.SplitAfter(rest, "/")
	}

	repo = &RepoURL{scheme, user, password, hostname, port, path}

	return
}

// ParseRepoURI is like ParseRepoURI, execpt also returns an error for legacy compatibility.
func ParseRepoURI(s string) (repo *RepoURL, err error) {
	return ParseRepoURL(s), nil
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
