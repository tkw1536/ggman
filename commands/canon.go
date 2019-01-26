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
	la := len(parsed.Args)
	if la > 2 || la == 0 {
		err = stringCanonTakesOneOrTwoArguments
		retval = ErrorSpecificParseArgs
		return
	}

	var lines []repos.CanLine
	var e error

	if la == 2 {
		// if we have two argument, use the specific specification given
		lines = append(lines, repos.CanLine{Pattern: "", Canonical: parsed.Args[1]})
	} else {
		// else read the canon file
		lines, e = getCanonOrPanic()
		if e != nil {
			err = stringInvalidCanfile
			retval = ErrorMissingConfig
			return
		}

	}

	// print the canonical url or error
	return printCanonOrError(lines, parsed.Args[0])
}

func printCanonOrError(lines []repos.CanLine, repo string) (retval int, err string) {
	// parse the repo uri
	uri, e := repos.NewRepoURI(repo)
	if e != nil {
		err = stringUnparsedRepoName
		retval = ErrorInvalidRepo
		return
	}

	// get the canonical one based on the canfile
	canonical := uri.CanonicalWith(lines)
	fmt.Println(canonical)

	return
}
