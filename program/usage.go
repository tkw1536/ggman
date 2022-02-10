package program

import (
	"fmt"
	"runtime"
	"time"

	"github.com/alessio/shellescape"
	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/program/usagefmt"
)

type Info struct {
	BuildVersion string
	BuildTime    time.Time

	MainName    string
	Description string
}

// FmtVersion formats version information to be shown to a human info
func (info Info) FmtVersion() string {
	return fmt.Sprintf("%s version %s, built %s, using %s", info.MainName, info.BuildVersion, info.BuildTime, runtime.Version())
}

// MainUsage returns a help page about ggman
func (p Program[Runtime, Parameters, Requirements]) MainUsage() usagefmt.Page {
	commands := append(p.Commands(), p.Aliases()...)

	return usagefmt.Page{
		MainName:    p.Info.MainName,
		MainOpts:    p.globalOptions(),
		Description: p.Info.Description,

		SubCommands: commands,
	}
}

// CommandUsage generates the usage information about a specific command
func (p Program[Runtime, Parameters, Requirements]) CommandUsage(cmdargs CommandArguments[Runtime, Parameters, Requirements]) usagefmt.Page {
	Description := cmdargs.Description

	return usagefmt.Page{
		MainName: p.Info.MainName,
		MainOpts: p.globalOptionsFor(Description.Requirements),

		Description: Description.Description,

		SubName: cmdargs.Arguments.Command,
		SubOpts: usagefmt.MakeOpts(cmdargs.Parser),

		MetaName: Description.PosArgName,
		MetaMin:  Description.PosArgsMin,
		MetaMax:  Description.PosArgsMax,

		Usage: Description.PosArgDescription,
	}
}

// AliasPage returns a usage page for the provided alias
func (p Program[Runtime, Parameters, Requirements]) AliasUsage(cmdargs CommandArguments[Runtime, Parameters, Requirements], alias Alias) usagefmt.Page {
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
		MainOpts: p.globalOptionsFor(cmdargs.Description.Requirements),

		Description: description,

		SubName: alias.Name,
		SubOpts: nil,

		MetaName: "ARG",
		MetaMin:  0,
		MetaMax:  -1,

		Usage: fmt.Sprintf("Arguments to pass after %s.", exCmd),
	}
}

var universalOpts = usagefmt.MakeOpts(flags.NewParser(&Universals{}, flags.None))

// globalOptions returns all global options
func (p Program[Runtime, Parameters, Requirements]) globalOptions() (opts []usagefmt.Opt) {
	opts = append(opts, universalOpts...)
	opts = append(opts, p.flagOptions()...)
	return
}

// globalOptionsFor returns global options for the provided requirement
func (p Program[Runtime, Parameters, Requirements]) globalOptionsFor(r Requirements) (opts []usagefmt.Opt) {
	flags := p.flagOptions()

	// filter options to be those that are allowed
	n := 0
	for _, opt := range flags {
		if !r.AllowsOption(opt) {
			continue
		}
		flags[n] = opt
		n++
	}
	flags = flags[:n]

	// add both universals and then locals
	opts = append(opts, universalOpts...)
	opts = append(opts, flags...)
	return
}

// flagOptions returns the options something something something
func (p Program[Runtime, Parameters, Requirements]) flagOptions() []usagefmt.Opt {
	return usagefmt.MakeOpts(flags.NewParser(&Flags{}, flags.None))
}
