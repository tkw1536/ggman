package program

import (
	"github.com/tkw1536/ggman/src/args"
)

// SubCommand represents a command that can be run with the program
type SubCommand func(args *args.GGArgs) (retval int, err string)

// Program represents a main program
type Program struct {
	commands   map[string]SubCommand
	wrapLength int
}
