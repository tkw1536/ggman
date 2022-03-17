package mockenv

import (
	"fmt"

	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/goprogram/stream"
)

// MappedPlumbing is a version of a git.Plumbing that runs all URLs via a translation mapping.
// It otherwise wraps an underlying plumbing that is used for all operations.
//
// It should be used in git operations that would otherwise require a real remote repository.
type MappedPlumbing struct {
	// Plumbing is the underlying plumbing used by this MappedPlumbing.
	git.Plumbing

	// URLMap is the map that contains the mapping that URLs are mapped via.
	// Keys correspond to URLs passed to this plumbing.
	// Values take the form of url passed to the underlying plumbing.
	URLMap map[string]string
}

// Forward maps a URL passed to this plumbing into a URL to the underlying plumbing.
// When the url does not exist in the mapping, calls panic().
func (mp MappedPlumbing) Forward(url string) string {
	translated, hasURL := mp.URLMap[url]
	if !hasURL {
		panic(fmt.Sprintf("MappedPlumbing: %q has no forward mapping", url))
	}
	return translated
}

// Backward maps a URL passed to the underlying plumbing into a URL to this plumbing.
// When the url does not exist in the mapping, calls panic().
func (mp MappedPlumbing) Backward(url string) string {
	for translated, u := range mp.URLMap {
		if u == url {
			return translated
		}
	}
	panic(fmt.Sprintf("MappedPlumbing: %q has no reverse mapping", url))
}

// Clone translates remoteURI and calls Clone on the underlying Plumbing.
func (mp MappedPlumbing) Clone(stream stream.IOStream, remoteURI, clonePath string, extraargs ...string) error {
	return mp.Plumbing.Clone(stream, mp.Forward(remoteURI), clonePath, extraargs...)
}

// GetRemotes calls GetRemotes() on the underlying Plumbing and translates the returned URLs.
func (mp MappedPlumbing) GetRemotes(clonePath string, repoObject any) (remotes map[string][]string, err error) {
	remotes, err = mp.Plumbing.GetRemotes(clonePath, repoObject)
	for k := range remotes {
		for i := range remotes[k] {
			remotes[k][i] = mp.Backward(remotes[k][i])
		}
	}
	return
}

// GetCanonicalRemote calls GetCanonicalRemote() on the underlying Plumbing and translates all returned urls.
func (mp MappedPlumbing) GetCanonicalRemote(clonePath string, repoObject any) (name string, urls []string, err error) {
	name, urls, err = mp.Plumbing.GetCanonicalRemote(clonePath, repoObject)
	for i := range urls {
		urls[i] = mp.Backward(urls[i])
	}
	return
}

// SetRemoteURLs translates urls and calls SetRemoteURLs() on the underlying Plumbing.
func (mp MappedPlumbing) SetRemoteURLs(clonePath string, repoObject any, name string, urls []string) (err error) {
	for i := range urls {
		urls[i] = mp.Forward(urls[i])
	}
	return mp.Plumbing.SetRemoteURLs(clonePath, repoObject, name, urls)
}

func init() {
	var _ git.Plumbing = (*MappedPlumbing)(nil)
}
