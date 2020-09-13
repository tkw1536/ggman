package env

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

// TODO: Test me

// Variables represents the values of specific environment variables.
// Unset variables are represented as the empty string.
type Variables struct {
	// HOME is the path to the users' home directory
	// This is typically stored in the 'HOME' variable on unix-like systems
	HOME string

	// PATH is the value of the 'PATH' environment variable
	PATH string

	// GGROOT is the value of the 'GGROOT' environment variable
	GGROOT string

	// CANFILE is the value of the 'GGMAN_CANFILE' environment variable
	CANFILE string
}

// ReadVariables reads a new variables instances from the environment
func ReadVariables() (v Variables) {
	v.HOME, _ = homedir.Dir() // errors result in an empty home value
	v.CANFILE = os.Getenv("GGMAN_CANFILE")
	v.GGROOT = os.Getenv("GGROOT")
	v.PATH = os.Getenv("PATH")
	return
}
