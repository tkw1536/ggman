package program

import (
	"fmt"
	"runtime"
	"time"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/ggman/program/usagefmt"
)

// Info contains meta-information about the current program
type Info struct {
	BuildVersion string
	BuildTime    time.Time

	MainName    string // Name of the main executable of the program
	Description string // Description of the program
}

// FmtVersion formats version information about the current version
// It returns a string that should be presented to users.
func (info Info) FmtVersion() string {
	return fmt.Sprintf("%s version %s, built %s, using %s", info.MainName, info.BuildVersion, info.BuildTime, runtime.Version())
}

// MainUsage returns a help page about ggman
func (p Program[E, P, F, R]) MainUsage() usagefmt.Page {
	commands := append(p.Commands(), p.Aliases()...)

	return usagefmt.Page{
		MainName:    p.Info.MainName,
		MainOpts:    globalOptions[F](),
		Description: p.Info.Description,

		SubCommands: commands,
	}
}

// CommandUsage generates the usage information about a specific command
func (p Program[E, P, F, R]) CommandUsage(context Context[E, P, F, R]) usagefmt.Page {
	Description := context.Description

	return usagefmt.Page{
		MainName: p.Info.MainName,
		MainOpts: globalOptionsFor[F](Description.Requirements),

		Description: Description.Description,

		SubName: context.Args.Command,
		SubOpts: usagefmt.MakeOpts(context.Parser),

		MetaName: Description.PosArgName,
		MetaMin:  Description.PosArgsMin,
		MetaMax:  Description.PosArgsMax,

		Usage: Description.PosArgDescription,
	}
}

// AliasPage returns a usage page for the provided alias
func (p Program[E, P, F, R]) AliasUsage(context Context[E, P, F, R], alias Alias) usagefmt.Page {
	exCmd := "`" + shellescape.QuoteCommand(append([]string{p.Info.MainName}, alias.Expansion()...)) + "`"
	helpCmd := "`" + shellescape.QuoteCommand([]string{p.Info.MainName, alias.Command, "--help"}) + "`"
	name := shellescape.Quote(alias.Command)

	var description string
	if alias.Description != "" {
		description = alias.Description + "\n\n"
	}
	description += fmt.Sprintf("Alias for %s. See %s for detailed help page about %s. ", exCmd, helpCmd, name)

	return usagefmt.Page{
		MainName: p.Info.MainName,
		MainOpts: globalOptionsFor[F](context.Description.Requirements),

		Description: description,

		SubName: alias.Name,
		SubOpts: nil,

		MetaName: "ARG",
		MetaMin:  0,
		MetaMax:  -1,

		Usage: fmt.Sprintf("Arguments to pass after %s.", exCmd),
	}
}
