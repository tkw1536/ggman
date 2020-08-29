package repos

import (
	"strings"

	"github.com/tkw1536/ggman/util"
)

// Components gets the components of a URL
func (url URL) Components() []string {

	// First split the path into components split by '/'.
	// and remove a '.git' from the last part.
	parts := util.RemoveEmpty(strings.Split(url.Path, "/"))
	lastPart := len(parts) - 1
	if lastPart >= 0 {
		parts[lastPart] = strings.TrimSuffix(parts[lastPart], ".git")

		// if we had a '/' before the .git, remove it
		if parts[lastPart] == "" {
			parts = parts[:lastPart]
		}
	}

	// Now prepend the hostname and user (unless it is 'git' or missing)
	components := make([]string, 0, 2+len(parts))
	if url.User != "" && url.User != "git" {
		components = append(components, url.HostName, url.User)
	} else {
		components = append(components, url.HostName)
	}
	return append(components, parts...)
}
