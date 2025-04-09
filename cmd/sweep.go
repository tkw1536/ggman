package cmd

//spellchecker:words github ggman internal walker goprogram exit
import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/internal/walker"
	"github.com/tkw1536/goprogram/exit"
)

//spellchecker:words GGROOT nolint wrapcheck

// Sweep is the 'ggman sweep' command.
//
// The sweep command can be used to identify non-git directories within the GGROOT directory which are empty, or contain only subdirectories which are empty recursively.
// Such directories are left behind after running the 'ggman relocate' command, or after manually deleting repositories.
// The command takes no arguments, and produces them in an order such that they can be passed to 'rmdir' and be deleted.
var Sweep ggman.Command = sweep{}

type sweep struct{}

func (sweep) Description() ggman.Description {
	return ggman.Description{
		Command:     "sweep",
		Description: "find empty folders in the root folder",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

var errScanEmpty = exit.Error{
	Message:  "error scanning for empty directories",
	ExitCode: exit.ExitGeneric,
}

func (sweep) Run(context ggman.Context) error {
	results, err := walker.Sweep(func(path string, root walker.FS, depth int) (stop bool) {
		return context.Environment.Git.IsRepository(path)
	}, walker.Params{
		Root: walker.NewRealFS(context.Environment.Root, false),
	})
	if err != nil {
		return errScanEmpty.WrapError(err) //nolint:wrapcheck
	}

	for _, r := range results {
		if _, err := context.Println(r); err != nil {
			return ggman.ErrGenericOutput.WrapError(err) //nolint:wrapcheck
		}
	}
	return nil
}
