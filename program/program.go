package program

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/util"
)

// Program represents a main program
type Program struct {
	commands   map[string]SubCommand
	options    map[string]*SubOptions
	wrapLength int
}

// NewProgram makes a new program
func NewProgram() *Program {
	cmds := make(map[string]SubCommand)
	opts := make(map[string]*SubOptions)
	return &Program{commands: cmds, options: opts, wrapLength: 80}
}

// Register registers a new sub-command with this program
func (pgrm *Program) Register(name string, sub SubCommand, opts *SubOptions) {
	pgrm.commands[name] = sub
	pgrm.options[name] = opts
}

// Run runs this program
func (pgrm *Program) Run(argv []string) (retval int, err string) {
	defer func() {
		err = util.WrapStringPreserveJ(err, pgrm.wrapLength)
	}()

	// parse the arguments
	parsed, err := ParseArgs(argv)
	if err != "" {
		retval = constants.ErrorGeneralParsArgs
		return
	}

	// for help, print help
	if parsed.Help {
		pgrm.Print(fmt.Sprintf(constants.StringUsage, constants.BuildVersion, constants.BuildTime, pgrm.knownCommands()))

		retval = 0
		err = ""
		return
	}

	if parsed.Version {
		pgrm.Print(fmt.Sprintf(constants.StringVersion, constants.BuildVersion, constants.BuildTime))

		retval = 0
		err = ""
		return
	}

	// extract the command and options
	cmd, ok := pgrm.commands[parsed.Command]
	opts, ok2 := pgrm.options[parsed.Command]

	// if we did not find both the command and options, something went wrong
	if !(ok && ok2) {
		err = fmt.Sprintf(constants.StringUnknownCommand, pgrm.knownCommands())
		retval = constants.ErrorUnknownCommand
		return
	}

	// prepare to parse options
	var runtime *SubRuntime
	var shouldContinue bool

	// apply the options to the argumen ts
	runtime, shouldContinue, retval, err = opts.Apply(pgrm, parsed)

	// if we do not continue, we should abort here
	if !shouldContinue {
		return
	}

	// finally run the command
	retval, err = cmd(runtime)

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

// Print prints output to the command line
func (pgrm *Program) Print(s string) {
	fmt.Println(util.WrapStringPreserveJ(s, pgrm.wrapLength))
}
