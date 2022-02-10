package program

import (
	"fmt"
	"runtime"
	"time"

	"github.com/alessio/shellescape"
	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/internal/text"
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
func (p Program[Runtime, Requirements]) MainUsage() usagefmt.Page {
	commands := append(p.Commands(), p.Aliases()...)

	return usagefmt.Page{
		MainName:    p.Info.MainName,
		MainOpts:    GetMainOpts(nil),
		Description: p.Info.Description,

		SubCommands: commands,
	}
}

// CommandUsage generates the usage information about a specific command
func (p Program[Runtime, Requirements]) CommandUsage(cmdargs CommandArguments[Runtime, Requirements]) usagefmt.Page {
	opt := cmdargs.description

	return usagefmt.Page{
		MainName: p.Info.MainName,
		MainOpts: GetMainOpts(&opt),

		Description: opt.Description,

		SubName: cmdargs.Arguments.Command,
		SubOpts: usagefmt.MakeOpts(cmdargs.parser),

		MetaName: opt.PosArgName,
		MetaMin:  opt.PosArgsMin,
		MetaMax:  opt.PosArgsMax,

		Usage: opt.PosArgDescription,
	}
}

// AliasPage returns a usage page for the provided alias
func (p Program[Runtime, Requirements]) AliasUsage(cmdargs CommandArguments[Runtime, Requirements], alias Alias) usagefmt.Page {
	opt := cmdargs.description

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

// GetMainOpts returns a list of global options for the provided command
func GetMainOpts(opt *Description) (opts []usagefmt.Opt) {

	// generate the main options by parsing the fake 'Arguments' struct.
	// return immediatly if global options only were requested
	opts = usagefmt.MakeOpts(flags.NewParser(&Arguments{}, flags.None))
	if opt == nil {
		return opts
	}

	n := 0
	for _, arg := range opts {
		// when the environment does not allow a filter, we only allow non-filter options!
		if !opt.Environment.AllowsFilter && !text.SliceContainsAny(arg.Long(), argumentsGeneralOptions...) {
			continue
		}
		opts[n] = arg
		n++
	}
	return opts[:n]
}
