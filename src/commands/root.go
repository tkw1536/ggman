package commands

import (
	"github.com/tkw1536/ggman/src/program"
)

// RootCommand is the entry point for the clone command
func RootCommand(runtime *program.SubRuntime) (retval int, err string) {

	// and echo out the root directory
	println(runtime.Root)

	// and exit
	return
}
