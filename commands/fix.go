package commands

import (
	"fmt"
	"os"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/repos"
)

// FixCommand is the entry point for the fix command
func FixCommand(parsed *GGArgs) (retval int, err string) {

	// read the --simulate flag
	simulateFlag, ie := parsed.ParseSingleFlag("--simulate")
	if ie {
		err = stringFixArguments
		retval = constants.ErrorSpecificParseArgs
		return
	}

	// get the canfile
	lines, e := getCanonOrPanic()
	if e != nil {
		err = stringInvalidCanfile
		retval = constants.ErrorMissingConfig
		return
	}

	// get the root directory or panic
	root, e := getRootOrPanic()
	if e != nil {
		err = stringUnableParseRootDirectory
		retval = constants.ErrorMissingConfig
		return
	}

	// find all the repos
	rs := repos.Repos(root, parsed.Pattern)
	hasError := false

	var msg string

	// and fix them all
	for _, repo := range rs {
		if simulateFlag {
			msg = fmt.Sprintf("Simulate fixing remote of %q", repo)
		} else {
			msg = fmt.Sprintf("Fixing remote of %q", repo)
		}
		if e := repos.FixRemote(repo, simulateFlag, msg, lines); e != nil {
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