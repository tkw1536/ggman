package program

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/internal/text"
	"github.com/tkw1536/ggman/program/usagefmt"
)

// UsagePage returns a help page about ggman
func (p Program) UsagePage() usagefmt.Page {
	text := "ggman manages local git repositories.\n\n"
	text += fmt.Sprintf("ggman version %s\n", constants.BuildVersion)
	text += "ggman is licensed under the terms of the MIT License.\nUse 'ggman license' to view licensing information."

	commands := append(p.Commands(), p.Aliases()...)

	return usagefmt.Page{
		MainName:    "ggman",
		MainOpts:    GetMainOpts(nil),
		Description: text,

		SubCommands: commands,
	}
}

// UsagePage generates a help page about this ggman subcommand
func (cmdargs CommandArguments) UsagePage() usagefmt.Page {
	opt := cmdargs.description

	return usagefmt.Page{
		MainName: "ggman",
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
