package program

import (
	"fmt"
	"sort"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/ggman/internal/usagefmt"
)

// Alias represents an alias for a different command
type Alias struct {
	Name string // the name of this alias

	Command string   // command to invoke
	Args    []string // arguments to pass to command

	Description string // Description of the alias
}

// Invoke returns the new command name and arguments when alias in invoked with the provided arguments
func (a Alias) Invoke(args []string) (command string, arguments []string) {
	// setup command
	command = a.Command

	// setup arguments
	arguments = make([]string, 0, len(a.Args)+len(args))
	arguments = append(arguments, a.Args...)
	arguments = append(arguments, args...)

	return
}

func (a Alias) Expansion() []string {
	return append([]string{a.Command}, a.Args...)
}

// RegisterAlias registers an alias.
//
// An alias must be of non-zero length.
// If an alias is of length zero, RegisterAlias calls panic().
// If an alias already exists, RegisterAlias calls panic().
//
// Aliases must not contain global flags; execution of the alias will fail at runtime.
// Aliases are not expanded recursively, meaning one alias may not refer to itself or another.
// An alias always takes precedence over a command with the same name.
func (p *Program) RegisterAlias(alias Alias) {
	if p.aliases == nil {
		p.aliases = make(map[string]Alias)
	}

	name := alias.Name
	if _, ok := p.aliases[name]; ok {
		panic("RegisterAlias(): Alias already registered")
	}

	p.aliases[name] = alias
}

func (p Program) Aliases() []string {
	aliases := make([]string, 0, len(p.aliases))
	for alias := range p.aliases {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	return aliases
}

// AliasPage returns a usage page for an alias
func (cmdargs CommandArguments) AliasPage(alias Alias) usagefmt.Page {
	opt := cmdargs.description

	exCmd := "`" + shellescape.QuoteCommand(append([]string{"ggman"}, alias.Expansion()...)) + "`"
	helpCmd := "`" + shellescape.QuoteCommand([]string{"ggman", alias.Command, "--help"}) + "`"
	name := shellescape.Quote(alias.Command)

	var description string
	if alias.Description != "" {
		description = alias.Description + "\n\n"
	}
	description += fmt.Sprintf("Alias for %s. See %s for detailed help page about %s. ", exCmd, helpCmd, name)

	return usagefmt.Page{
		MainName: "ggman",
		MainOpts: GetMainOpts(&opt),

		Description: description,

		SubName: alias.Name,
		SubOpts: nil,

		MetaName: "ARG",
		MetaMin:  0,
		MetaMax:  -1,

		Usage: fmt.Sprintf("Arguments to pass after %s.", exCmd),
	}
}
