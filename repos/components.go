package repos

import (
	"regexp"
	"strings"

	"github.com/tkw1536/ggman/util"
)

var reSlash = regexp.MustCompile("/+")

// Components gets the Components of a Repo URI
func (rURI *RepoURI) Components() (parts []string) {

	// normalize the path
	path := (*rURI).Path
	path = util.TrimPrefixWhile(util.TrimSuffixWhile(path, "/"), "/")
	path = strings.TrimSuffix(path, ".git")
	path = util.TrimSuffixWhile(path, "/")
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
