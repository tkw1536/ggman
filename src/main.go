package src

import (
	"github.com/tkw1536/ggman/src/commands"
	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/program"
)

const tWL = 80

// Main is the main entry point for the program
func Main(argv []string) (retval int, err string) {

	// make a new program
	ggman := program.NewProgram()

	// register all the commands
	ggman.Register("root", commands.RootCommand, &program.SubOptions{NeedsRoot: true})

	ggman.Register("ls", commands.LSCommand, &program.SubOptions{ForArgument: program.OptionalFor, Flag: "--exit-code", FlagDescription: constants.StringExitFlagUsage, NeedsRoot: true})
	ggman.Register("lsr", commands.LSRCommand, &program.SubOptions{ForArgument: program.OptionalFor, Flag: "--canonical", FlagDescription: constants.StringCanonicalFlagUsage, NeedsRoot: true})

	ggman.Register("where", commands.WhereCommand, &program.SubOptions{MinArgs: 1, MaxArgs: 1, Metavar: "REPO", UsageDescription: constants.StringWhereRepoUsage, NeedsRoot: true})

	ggman.Register("canon", commands.CanonCommand, &program.SubOptions{MinArgs: 1, MaxArgs: 2, UsageDescription: constants.StringCanonUsage})
	ggman.Register("comps", commands.CompsCommand, &program.SubOptions{MinArgs: 1, MaxArgs: 1, Metavar: "URI", UsageDescription: constants.StringCompsURIUsage})

	ggman.Register("fetch", commands.FetchCommand, &program.SubOptions{ForArgument: program.OptionalFor, NeedsRoot: true})
	ggman.Register("pull", commands.PullCommand, &program.SubOptions{ForArgument: program.OptionalFor, NeedsRoot: true})

	ggman.Register("fix", commands.FixCommand, &program.SubOptions{ForArgument: program.OptionalFor, Flag: "--simulate", FlagDescription: constants.StringSimulateFlagUsage, NeedsRoot: true, NeedsCANFILE: true})

	ggman.Register("clone", commands.CloneCommand, &program.SubOptions{MinArgs: 1, MaxArgs: -1, Metavar: "ARG", UsageDescription: constants.StringCloneURIUsage, NeedsRoot: true, NeedsCANFILE: true})
	ggman.Register("link", commands.LinkCommand, &program.SubOptions{MinArgs: 1, MaxArgs: 1, Metavar: "PATH", UsageDescription: constants.StringLinkPathUsage, NeedsRoot: true})

	ggman.Register("license", commands.LicenseCommand, &program.SubOptions{})

	ggman.Register("here", commands.HereCommand, &program.SubOptions{NeedsRoot: true, Flag: "--tree", FlagDescription: constants.StringTreeFlagUsage})

	ggman.Register("web", commands.WebCommand, &program.SubOptions{NeedsRoot: true, Flag: "--tree", MinArgs: 0, MaxArgs: 1, Metavar: "BASE", FlagDescription: constants.StringTreeFlagUsage})
	ggman.Register("url", commands.URLCommand, &program.SubOptions{NeedsRoot: true, Flag: "--tree", MinArgs: 0, MaxArgs: 1, Metavar: "BASE", FlagDescription: constants.StringTreeFlagUsage})

	// and run it
	return ggman.Run(argv)
}
