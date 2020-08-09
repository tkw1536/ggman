package commands

import (
	"fmt"
	"os"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/repos"
)

// FixCommand is the entry point for the fix command
func FixCommand(runtime *program.SubRuntime) (retval int, err string) {

	// read runtime info
	simulateFlag := runtime.Flag
	lines := runtime.Canfile
	root := runtime.Root

	// find all the repos
	rs := repos.Repos(root, runtime.For)
	hasError := false

	// and fix them all
	for _, repo := range rs {

		// we might have to print a log message
		var initialMessage string
		didPrintInitialMessage := false
		if simulateFlag {
			initialMessage = fmt.Sprintf("Simulate fixing remote of %q", repo)
		} else {
			initialMessage = fmt.Sprintf("Fixing remote of %q", repo)
		}

		if e := git.Default.UpdateRemotes(repo, func(url, remoteName string) (string, error) {

			// print a log message if we haven't already
			if !didPrintInitialMessage {
				didPrintInitialMessage = true
				fmt.Println(initialMessage)
			}

			// compute the new canonical url
			canon := repos.ParseRepoURL(url).CanonicalWith(lines)

			fmt.Printf("Updating %s: %s -> %s\n", remoteName, url, canon)

			// either return the canonical url, or (if we're simulating) the old url
			if simulateFlag {
				return url, nil
			}

			return canon, nil
		}); e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			hasError = true
		}
	}

	// if we had an error, indicate that to the user
	if hasError {
		retval = constants.ErrorCodeCustom
	}

	// and finish
	return
}
