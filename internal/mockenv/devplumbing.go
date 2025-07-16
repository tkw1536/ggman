//spellchecker:words mockenv
package mockenv

//spellchecker:words strconv ggman pkglib stream
import (
	"fmt"
	"io"
	"strconv"

	"go.tkw01536.de/ggman/git"
	"go.tkw01536.de/pkglib/stream"
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

// stream returns a mapped version of stream to be used.
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
		panic("DevPlumbing: " + strconv.Quote(url) + " has no forward mapping")
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
func (dp DevPlumbing) Clone(stream stream.IOStream, remoteURI, clonePath string, extraArgs ...string) error {
	err := dp.Plumbing.Clone(dp.stream(stream), dp.Forward(remoteURI), clonePath, extraArgs...)
	if err != nil {
		return fmt.Errorf("%q: failed to clone: %w", remoteURI, err)
	}
	return nil
}

// Fetch called Fetch on the underlying Plumbing.
func (dp DevPlumbing) Fetch(stream stream.IOStream, clonePath string, cache any) error {
	err := dp.Plumbing.Fetch(dp.stream(stream), clonePath, cache)
	if err != nil {
		return fmt.Errorf("%q: failed to fetch: %w", clonePath, err)
	}
	return nil
}

// Fetch calls Pull on the underlying Plumbing.
func (dp DevPlumbing) Pull(stream stream.IOStream, clonePath string, cache any) error {
	err := dp.Plumbing.Pull(dp.stream(stream), clonePath, cache)
	if err != nil {
		return fmt.Errorf("%q: failed to pull: %w", clonePath, err)
	}
	return nil
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
func (dp DevPlumbing) SetRemoteURLs(clonePath string, repoObject any, name string, urls []string) error {
	for i := range urls {
		urls[i] = dp.Forward(urls[i])
	}

	err := dp.Plumbing.SetRemoteURLs(clonePath, repoObject, name, urls)
	if err != nil {
		return fmt.Errorf("%q: failed to set remote URLs: %w", clonePath, err)
	}
	return nil
}

func init() {
	var _ git.Plumbing = (*DevPlumbing)(nil)
}
