package commands

import (
	"fmt"
	"os"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/gitwrap"
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

	var msg string

	// and fix them all
	for _, repo := range rs {
		if simulateFlag {
			msg = fmt.Sprintf("Simulate fixing remote of %q", repo)
		} else {
			msg = fmt.Sprintf("Fixing remote of %q", repo)
		}
		if e := gitwrap.Implementation.FixRemotes(repo, simulateFlag, msg, lines); e != nil {
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
