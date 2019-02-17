package repos

import (
	"errors"
	"fmt"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

// GetRemote gets the remote of a repository
func GetRemote(filepath string) (uri string, err error) {
	// open filepath or error out
	repo, err := git.PlainOpen(filepath)
	if err != nil {
		return
	}

	// get the current head
	remotes, err := repo.Remotes()
	if err != nil {
		return
	}

	// get the canonical remote (commonly named 'origin')
	canon, err := getCanonicalRemote(remotes, "origin")
	if err != nil {
		return
	}

	// and return it's fetching url
	uri = canon.Config().URLs[0]
	return
}

// FixRemote updates the remote of a repository with a given CanLine array
func FixRemote(filepath string, simulate bool, initialLogLine string, lines []CanLine) (err error) {
	// if we did not print anything, we need to still print something
	didPrint := false
	defer (func() {
		if !didPrint && err != nil {
			fmt.Println(initialLogLine)
		}
	})()

	// open filepath or error out
	r, err := git.PlainOpen(filepath)
	if err != nil {
		return
	}

	// get the current head
	remotes, err := r.Remotes()
	if err != nil {
		return
	}

	// get the canonical remote (commonly named 'origin')
	canon, err := getCanonicalRemote(remotes, "origin")
	if err != nil {
		return
	}

	// fetch the current configuration
	cfg, err := r.Storer.Config()
	if err != nil {
		return
	}

	// update the urls
	didPrint, cfg.Remotes[canon.Config().Name].URLs = fixURLs(canon.Config(), initialLogLine, lines)

	// and store it again
	if !simulate {
		err = r.Storer.SetConfig(cfg)
	}
	return
}

// fixURLs fixes the urls of a remote
func fixURLs(config *config.RemoteConfig, initialLogLine string, lines []CanLine) (didPrint bool, fixed []string) {
	didPrint = false

	for _, url := range config.URLs {
		current, err := NewRepoURI(url)
		if err != nil {
			continue
		}
		canon := current.CanonicalWith(lines)
		if canon != url {
			if !didPrint {
				fmt.Println(initialLogLine)
				didPrint = true
			}
			fmt.Printf("Updating %s: %s -> %s\n", config.Name, url, canon)
		}
		fixed = append(fixed, canon)
	}
	return
}

// gets the canonical remote of a repository
func getCanonicalRemote(remotes []*git.Remote, name string) (canonical *git.Remote, err error) {
	// find the main remote
	remoteCount := len(remotes)

	// no remotes => error
	if remoteCount == 0 {
		err = errors.New("No remotes found")
		return
	}

	// if we have a remote named 'origin' use it
	for _, r := range remotes {
		if r.Config().Name == name {
			canonical = r
		}
	}

	// else pick the first remote
	if canonical == nil {
		canonical = remotes[0]
	}

	return
}
