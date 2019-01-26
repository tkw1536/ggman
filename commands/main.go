package commands

import "os"

// Main is the main entry point for the program
func Main(args []string) (retval int, err string) {
	// parse the arguments
	parsed, err := ParseArgs(os.Args[1:])
	if err != "" {
		retval = ErrorGeneralParsArgs
		return
	}

	// run the appropriate command
	switch parsed.Command {
	case "root":
		retval, err = RootCommand(parsed)
	case "ls":
		retval, err = LSCommand(parsed)
	case "where":
		retval, err = WhereCommand(parsed)
	default:
		err = stringUnknownCommand
		retval = ErrorUnknownCommand
	}

	// and return
	return
}
