package program

import (
	"fmt"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/ggman/program/meta"
)

// MainUsage returns a help page about ggman
func (p Program[E, P, F, R]) MainUsage() meta.Meta {
	commands := append(p.Commands(), p.Aliases()...)

	return meta.Meta{
		Executable:  p.Info.Executable,
		GlobalFlags: globalOptions[F](),
		Description: p.Info.Description,

		Commands: commands,
	}
}

// CommandUsage generates the usage information about a specific command
func (p Program[E, P, F, R]) CommandUsage(context Context[E, P, F, R]) meta.Meta {
	return meta.Meta{
		Executable:  p.Info.Executable,
		GlobalFlags: globalFlagsFor[F](context.Description.Requirements),

		Description: context.Description.Description,

		Command:      context.Description.Command,
		CommandFlags: meta.AllFlags(context.commandParser),

		Positional: context.Description.Positional,
	}
}

// AliasPage returns a usage page for the provided alias
func (p Program[E, P, F, R]) AliasUsage(context Context[E, P, F, R], alias Alias) meta.Meta {
	exCmd := "`" + shellescape.QuoteCommand(append([]string{p.Info.Executable}, alias.Expansion()...)) + "`"
	helpCmd := "`" + shellescape.QuoteCommand([]string{p.Info.Executable, alias.Command, "--help"}) + "`"
	name := shellescape.Quote(alias.Command)

	var description string
	if alias.Description != "" {
		description = alias.Description + "\n\n"
	}
	description += fmt.Sprintf("Alias for %s. See %s for detailed help page about %s. ", exCmd, helpCmd, name)

	return meta.Meta{
		Executable:  p.Info.Executable,
		GlobalFlags: globalFlagsFor[F](context.Description.Requirements),

		Description: description,

		Command:      alias.Name,
		CommandFlags: nil,

		Positional: meta.Positional{
			Value:       "ARG",
			Description: fmt.Sprintf("Arguments to pass after %s.", exCmd),
			Min:         0,
			Max:         -1,
		},
	}
}
