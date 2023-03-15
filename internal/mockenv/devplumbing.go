package mockenv

import (
	"fmt"
	"io"

	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/pkglib/stream"
)

// DevPlumbing is a version of a git.Plumbing that runs all URLs via a translation mapping
// and can optionally silence standard error output.
// It otherwise wraps an underlying plumbing that is used for all operations.
//
// It is intended for use during for testing.
// It should be used in git operations that would otherwise require a real remote repository.
type DevPlumbing struct {
	// Plumbing is the underlying plumbing used by this DevPlumbing.
	git.Plumbing

	// when set to true, silence standard error output for any operations.
	SilenceStderr bool

	// URLMap is the map that contains the mapping that URLs are mapped via.
	// Keys correspond to URLs passed to this plumbing.
	// Values take the form of url passed to the underlying plumbing.
	URLMap map[string]string
}

// stream returns a mapped version of stream to be used
func (dp DevPlumbing) stream(stream stream.IOStream) stream.IOStream {
	if dp.SilenceStderr {
		stream.Stderr = io.Discard
	}
	return stream
}

// Forward maps a URL passed to this plumbing into a URL to the underlying plumbing.
// When the url does not exist in the mapping, calls panic().
func (dp DevPlumbing) Forward(url string) string {
	translated, hasURL := dp.URLMap[url]
	if !hasURL {
		panic(fmt.Sprintf("DevPlumbing: %q has no forward mapping", url))
	}
	return translated
}

// Backward maps a URL passed to the underlying plumbing into a URL to this plumbing.
// When the url does not exist in the mapping, calls panic().
func (dp DevPlumbing) Backward(url string) string {
	for translated, u := range dp.URLMap {
		if u == url {
			return translated
		}
	}
	panic(fmt.Sprintf("DevPlumbing: %q has no reverse mapping", url))
}

// Clone translates remoteURI and calls Clone on the underlying Plumbing.
func (dp DevPlumbing) Clone(stream stream.IOStream, remoteURI, clonePath string, extraargs ...string) error {
	return dp.Plumbing.Clone(dp.stream(stream), dp.Forward(remoteURI), clonePath, extraargs...)
}

// Fetch called Fetch on the underlying Plumbing
func (dp DevPlumbing) Fetch(stream stream.IOStream, clonePath string, cache any) (err error) {
	return dp.Plumbing.Fetch(dp.stream(stream), clonePath, cache)
}

// Fetch calls Pull on the underlying Plumbing
func (dp DevPlumbing) Pull(stream stream.IOStream, clonePath string, cache any) (err error) {
	return dp.Plumbing.Pull(dp.stream(stream), clonePath, cache)
}

// GetRemotes calls GetRemotes() on the underlying Plumbing and translates the returned URLs.
func (dp DevPlumbing) GetRemotes(clonePath string, repoObject any) (remotes map[string][]string, err error) {
	remotes, err = dp.Plumbing.GetRemotes(clonePath, repoObject)
	for k := range remotes {
		for i := range remotes[k] {
			remotes[k][i] = dp.Backward(remotes[k][i])
		}
	}
	return
}

// GetCanonicalRemote calls GetCanonicalRemote() on the underlying Plumbing and translates all returned urls.
func (dp DevPlumbing) GetCanonicalRemote(clonePath string, repoObject any) (name string, urls []string, err error) {
	name, urls, err = dp.Plumbing.GetCanonicalRemote(clonePath, repoObject)
	for i := range urls {
		urls[i] = dp.Backward(urls[i])
	}
	return
}

// SetRemoteURLs translates urls and calls SetRemoteURLs() on the underlying Plumbing.
func (dp DevPlumbing) SetRemoteURLs(clonePath string, repoObject any, name string, urls []string) (err error) {
	for i := range urls {
		urls[i] = dp.Forward(urls[i])
	}
	return dp.Plumbing.SetRemoteURLs(clonePath, repoObject, name, urls)
}

func init() {
	var _ git.Plumbing = (*DevPlumbing)(nil)
}
