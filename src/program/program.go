package program

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/utils"
)

// Program represents a main program
type Program struct {
	commands   map[string]SubCommand
	wrapLength int
}

// NewProgram makes a new program
func NewProgram() *Program {
	cmds := make(map[string]SubCommand)
	return &Program{commands: cmds, wrapLength: 80}
}

// Register registers a new sub-command with this program
func (pgrm *Program) Register(name string, sub SubCommand) {
	pgrm.commands[name] = sub
}

// Run runs this program
func (pgrm *Program) Run(argv []string) (retval int, err string) {
	defer func() {
		err = utils.WrapStringPreserveJ(err, pgrm.wrapLength)
	}()

	// parse the arguments
	parsed, err := ParseArgs(argv)
	if err != "" {
		retval = constants.ErrorGeneralParsArgs
		return
	}

	// for help, print help
	if parsed.Help {
		fmt.Println(utils.WrapStringPreserveJ(
			fmt.Sprintf(constants.StringUsage, pgrm.knownCommands()),
			pgrm.wrapLength))

		retval = 0
		err = ""
		return
	}

	// if we have a known command, run it
	if val, ok := pgrm.commands[parsed.Command]; ok {
		retval, err = val(parsed)

		// else throw an error
	} else {
		err = fmt.Sprintf(constants.StringUnknownCommand, pgrm.knownCommands())
		retval = constants.ErrorUnknownCommand
	}

	// and return
	return
}

func (pgrm *Program) knownCommands() string {
	keys := make([]string, 0, len(pgrm.commands))
	for k := range pgrm.commands {
		keys = append(keys, "'"+k+"'")
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}
