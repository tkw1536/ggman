package utils

import (
	"net/url"
	"regexp"
	"strings"
)

// Components takes a git-like url and splits it into a list of components
func Components(s string) (parts []string, err error) {

	// parse it into a url
	u, err := url.Parse(componentsNormURL(s))
	if err != nil {
		return
	}

	// extract hostname and path and normalize them
	hostname := componentsNormHostName(u.Hostname())
	path := componentsNormPath(u.Path)

	// extract the username
	var username string
	user := u.User
	if user != nil {
		username = strings.ToLower(user.Username())
	}

	parts = strings.Split(path, "/")

	// if we have a username that is not "git"
	// prepend it to the path
	if username != "" && username != "git" {
		parts = append([]string{hostname, username}, parts...)
	} else {
		parts = append([]string{hostname}, parts...)
	}

	return
}

// componentsNormURL norms a path to be used by te components function
func componentsNormURL(s string) string {
	// split into []string{protocol, rest} and fallback to unknown protocol
	splits := strings.SplitN(s, "://", 2)
	if len(splits) == 1 {
		splits = append([]string{"unknown"}, splits...)
	}

	return splits[0] + "://" + reColon.ReplaceAllString(splits[1], "/$1$2")
}

func componentsNormPath(s string) (norm string) {
	norm = TrimPrefixWhile(TrimSuffixWhile(strings.TrimSuffix(s, ".git"), "/"), "/")
	return reSlash.ReplaceAllString(norm, "/")
}

func componentsNormHostName(s string) string {
	// if we have an IPV6 address, remove the []s around it
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return strings.TrimSuffix(strings.TrimPrefix("[", s), "]")
	}

	// else return it as-is
	return s
}

var reSlash *regexp.Regexp // regular expression matching repeated slashes
var reColon *regexp.Regexp // regular expression matching colons in path that are not ports

func init() {
	reSlash = regexp.MustCompile("/+")
	reColon = regexp.MustCompile(":([^\\d][^/]*)(/|$)")
}
