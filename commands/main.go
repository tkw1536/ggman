package commands

import (
	"fmt"
	"os"

	"github.com/tkw1536/ggman/constants"
)

// Main is the main entry point for the program
func Main(args []string) (retval int, err string) {
	// parse the arguments
	parsed, err := ParseArgs(os.Args[1:])
	if err != "" {
		retval = constants.ErrorGeneralParsArgs
		return
	}

	// for help, print help
	if parsed.Help {
		fmt.Println(stringUsage)
		retval = 0
		err = ""
		return
	}

	// run the appropriate command
	switch parsed.Command {
	case "root":
		retval, err = RootCommand(parsed)
	case "ls":
		retval, err = LSCommand(parsed)
	case "lsr":
		retval, err = LSRCommand(parsed)
	case "where":
		retval, err = WhereCommand(parsed)
	case "canon":
		retval, err = CanonCommand(parsed)
	case "comps":
		retval, err = CompsCommand(parsed)
	case "fetch":
		retval, err = FetchCommand(parsed)
	case "pull":
		retval, err = PullCommand(parsed)
	case "fix":
		retval, err = FixCommand(parsed)
	case "clone":
		retval, err = CloneCommand(parsed)
	case "link":
		retval, err = LinkCommand(parsed)
	case "license":
		retval, err = LicenseCommand(parsed)
	default:
		err = stringUnknownCommand
		retval = constants.ErrorUnknownCommand
	}

	// and return
	return
}
