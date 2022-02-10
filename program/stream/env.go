package stream

import "os"

// NewEnvIOStream creates a new IOStream using the environment.
//
// The Stdin, Stdout and Stderr streams are used from the os package.
func NewEnvIOStream() IOStream {
	return IOStream{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		wrap:   ioDefaultWrap,
	}
}
