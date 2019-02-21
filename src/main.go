package src

import (
	"github.com/tkw1536/ggman/src/commands"
	"github.com/tkw1536/ggman/src/program"
)

const tWL = 80

// Main is the main entry point for the program
func Main(argv []string) (retval int, err string) {

	// make a new program
	ggman := program.NewProgram()

	// register all the commands
	ggman.Register("root", commands.RootCommand)
	ggman.Register("ls", commands.LSCommand)
	ggman.Register("lsr", commands.LSRCommand)
	ggman.Register("where", commands.WhereCommand)
	ggman.Register("canon", commands.CanonCommand)
	ggman.Register("comps", commands.CompsCommand)
	ggman.Register("fetch", commands.FetchCommand)
	ggman.Register("pull", commands.PullCommand)
	ggman.Register("fix", commands.FixCommand)
	ggman.Register("clone", commands.CloneCommand)
	ggman.Register("link", commands.LinkCommand)
	ggman.Register("license", commands.LicenseCommand)

	// and run it
	return ggman.Run(argv)
}
