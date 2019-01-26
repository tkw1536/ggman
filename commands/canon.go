package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/repos"
)

// CanonCommand is the entry point for the canon command
func CanonCommand(parsed *GGArgs) (retval int, err string) {
	// 'canon' takes no for
	if parsed.Pattern != "" {
		err = stringCanonNoFor
		retval = ErrorSpecificParseArgs
		return
	}

	// we accept one arguments
	if len(parsed.Args) != 1 {
		err = stringCanonTakesOneArgument
		retval = ErrorSpecificParseArgs
		return
	}

	// read the canon file
	lines, e := getCanonOrPanic()
	if e != nil {
		err = stringInvalidCanfile
		retval = ErrorMissingConfig
		return
	}

	// parse the repo uri
	uri, e := repos.NewRepoURI(parsed.Args[0])
	if e != nil {
		err = stringUnparsedRepoName
		retval = ErrorInvalidRepo
		return
	}

	// get the canonical one based on the canfile
	canonical := uri.CanonicalWith(lines)
	fmt.Println(canonical)

	// and finish
	return
}
