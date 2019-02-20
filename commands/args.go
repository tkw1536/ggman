package commands

import (
	"os"

	"github.com/tkw1536/ggman/repos"

	homedir "github.com/mitchellh/go-homedir"
)

// GGArgs represents the arguments passed to a gg command
// gg [for $pattern] $command [$args...]
type GGArgs struct {
	Command string
	Pattern string
	Help    bool
	Args    []string
}

const helpLongForm = "--help"
const helpShortForm = "-h"
const helpLiteralForm = "help"

const forLongForm = "--for"
const forShortForm = "-f"
const forLiteralForm = "for"

// ParseArgs parses arguments from the command line
func ParseArgs(args []string) (parsed *GGArgs, err string) {

	// if we have no arguments, that is an error
	count := len(args)
	if count == 0 {
		err = stringNeedOneArgument
		return
	}

	// if the first argument is help, everything else is ignroed
	head := args[0]

	// if the first argument is help, return help
	if head == helpLiteralForm || head == helpShortForm || head == helpLongForm {
		parsed = &GGArgs{"", "", true, args[1:]}
		return
	}

	if head == forLiteralForm || head == forShortForm || head == forLongForm {
		// gg for $pattern $command
		if count < 3 {
			err = stringNeedTwoAfterFor
			return
		}

		parsed = &GGArgs{args[2], args[1], false, args[3:]}
		return
	}

	// gg $pattern $command
	parsed = &GGArgs{args[0], "", false, args[1:]}
	return
}

// ParseSingleFlag parses a single optional flag
func (parsed *GGArgs) ParseSingleFlag(flag string) (value bool, err bool) {
	la := len(parsed.Args)

	// if we have too many arguments throw an error
	if la > 1 {
		err = true
		return
	} else if la == 1 && parsed.Args[0] != flag {
		err = true
		return
	}

	// and return the error
	value = la == 1
	err = false
	return
}

func getRootOrPanic() (value string, err error) {
	value = os.Getenv("GGROOT")
	if len(value) == 0 {
		value, err = homedir.Expand("~/Projects")
	}

	return
}

func getCanonOrPanic() (lines []repos.CanLine, err error) {
	return repos.ReadDefaultCanFile()
}
