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

	ggman.Register("ls", commands.LSCommand, &program.SubOptions{ForArgument: program.OptionalFor, Flag: "--exit-code", UsageDescription: constants.StringExitFlagUsage, NeedsRoot: true})
	ggman.Register("lsr", commands.LSRCommand, &program.SubOptions{ForArgument: program.OptionalFor, Flag: "--canonical", UsageDescription: constants.StringCanonicalFlagUsage, NeedsRoot: true})

	ggman.Register("where", commands.WhereCommand, &program.SubOptions{MinArgs: 1, MaxArgs: 1, Metavar: "REPO", UsageDescription: constants.StringWhereRepoUsage, NeedsRoot: true})

	ggman.Register("canon", commands.CanonCommand, &program.SubOptions{MinArgs: 1, MaxArgs: 2, UsageDescription: constants.StringCanonUsage})
	ggman.Register("comps", commands.CompsCommand, &program.SubOptions{MinArgs: 1, MaxArgs: 1, Metavar: "URI", UsageDescription: constants.StringCompsURIUsage})

	ggman.Register("fetch", commands.FetchCommand, &program.SubOptions{ForArgument: program.OptionalFor, NeedsRoot: true})
	ggman.Register("pull", commands.PullCommand, &program.SubOptions{ForArgument: program.OptionalFor, NeedsRoot: true})

	ggman.Register("fix", commands.FixCommand, &program.SubOptions{ForArgument: program.OptionalFor, Flag: "--simulate", UsageDescription: constants.StringSimulateFlagUsage, NeedsRoot: true, NeedsCANFILE: true})

	ggman.Register("clone", commands.CloneCommand, &program.SubOptions{MinArgs: 1, MaxArgs: -1, Metavar: "ARG", UsageDescription: constants.StringCloneURIUsage, NeedsRoot: true, NeedsCANFILE: true})
	ggman.Register("link", commands.LinkCommand, &program.SubOptions{MinArgs: 1, MaxArgs: 1, Metavar: "PATH", UsageDescription: constants.StringLinkPathUsage, NeedsRoot: true})

	ggman.Register("license", commands.LicenseCommand, &program.SubOptions{})

	// and run it
	return ggman.Run(argv)
}
