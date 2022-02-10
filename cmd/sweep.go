package cmd

import (
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/gg"
	"github.com/tkw1536/ggman/internal/walker"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/program/exit"
)

// Sweep is the 'ggman sweep' command.
//
// The sweep command can be used to identify non-git directories within the GGROOT directory which are empty, or contain only subdirectories which are empty recursively.
// Such directories are left behind after running the 'ggman relocate' command, or after manually deleting repositories.
// The command takes no arguments, and produces them in an order such that they can be passed to 'rmdir' and be deleted.
var Sweep program.Command = sweep{}

type sweep struct{}

func (sweep) BeforeRegister(program *program.Program) {}

func (sweep) Description() program.Description {
	return program.Description{
		Name:        "sweep",
		Description: "Find empty folders in the project folder. ",

		Environment: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (sweep) AfterParse() error {
	return nil
}

var errSweepErr = exit.Error{
	Message:  "Error scanning for empty directories: %s",
	ExitCode: exit.ExitGeneric,
}

func (sweep) Run(context program.Context) error {
	results, err := walker.Sweep(func(path string, root walker.FS, depth int) (stop bool) {
		return gg.C2E(context).Git.IsRepository(path)
	}, walker.Params{
		Root: walker.NewRealFS(gg.C2E(context).Root, false),
	})
	if err != nil {
		return errSweepErr.WithMessageF(err)
	}

	for _, r := range results {
		context.Println(r)
	}
	return nil
}
