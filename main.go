package main

import (
	"os"

	"github.com/tkw1536/ggman/commands"
)

func main() {
	// run the main command
	retval, err := commands.Main(os.Args[1:])

	// if we are not exiting with status code zero, print the error message
	// and then exit
	if retval != 0 {
		os.Stderr.WriteString(err + "\n")
		defer os.Exit(retval)
	}
}
