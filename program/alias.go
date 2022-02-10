package program

import (
	"sort"
)

// Alias represents an alias for a command.
//
// Expansion of an alias takes place at runtime.
// Aliases must not contain global flags; execution of the them will fail at runtime.
//
// Aliases are not expanded recursively, meaning one alias may not refer to itself or another.
// An alias always takes precedence over a command with the same name.
type Alias struct {
	// Name is the name of this alias
	Name string

	// Command to invoke along with arguments
	Command string
	Args    []string

	// Description for the usage page
	Description string
}

// Invoke returns command arguments that are to be used when invoking this alias
// args are additional arguments to pass to the command
func (a Alias) Invoke(args []string) (command string, arguments []string) {
	// setup command
	command = a.Command

	// setup arguments
	arguments = make([]string, 0, len(a.Args)+len(args))
	arguments = append(arguments, a.Args...)
	arguments = append(arguments, args...)

	return
}

// Expansion returns a slice representing the expansion of this alias.
func (a Alias) Expansion() []string {
	return append([]string{a.Command}, a.Args...)
}

// RegisterAlias registers a new alias.
// See also Alias.
//
// If an alias already exists, RegisterAlias calls panic().
func (p *Program[Runtime]) RegisterAlias(alias Alias) {
	if p.aliases == nil {
		p.aliases = make(map[string]Alias)
	}

	name := alias.Name
	if _, ok := p.aliases[name]; ok {
		panic("RegisterAlias(): Alias already registered")
	}

	p.aliases[name] = alias
}

// Aliases returns the names of all registered aliases.
// Aliases are returned in sorted order.
func (p Program[Runtime]) Aliases() []string {
	aliases := make([]string, 0, len(p.aliases))
	for alias := range p.aliases {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	return aliases
}
