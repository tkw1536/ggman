package cmd

//spellchecker:words path filepath github ggman internal dirs goprogram exit pkglib
import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/internal/dirs"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/fsx"
)

//spellchecker:words positionals nolint wrapcheck

// Link is the 'ggman link' command.
//
// The 'ggman link' symlinks the repository in the path passed as the first argument where it would have been cloned to inside 'ggman root'.
var Link ggman.Command = link{}

type link struct {
	Positionals struct {
		Path string `description:"path of repository to symlink" positional-arg-name:"PATH" required:"1-1"`
	} `positional-args:"true"`
}

func (link) Description() ggman.Description {
	return ggman.Description{
		Command:     "link",
		Description: "symlink a repository into the local repository structure",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

var errLinkDoesNotExist = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to link repository: can not open source repository",
}

var errLinkSamePath = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to link repository: link source and target are identical",
}

var errLinkAlreadyExists = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to link repository: another directory already exists in target location",
}

var errLinkUnknown = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unknown linking error",
}

func (l link) Run(context ggman.Context) error {
	// make sure that the path is absolute
	// to avoid relative symlinks
	from, e := context.Environment.Abs(l.Positionals.Path)
	if e != nil {
		return errLinkDoesNotExist
	}

	// open the source repository and get the remote
	r, e := context.Environment.Git.GetRemote(from, "")
	if e != nil {
		return errLinkDoesNotExist
	}

	// find the target path
	to, err := context.Environment.Local(env.ParseURL(r))
	if err != nil {
		return fmt.Errorf("%w: %w", errUnableLocalPath, err)
	}
	parentTo := filepath.Dir(to)

	// if it's the same path, we throw an error
	if from == to {
		return errLinkSamePath
	}

	// make sure it doesn't exist
	{
		exists, err := fsx.Exists(to)
		if err != nil {
			return fmt.Errorf("failed to check for existence: %w", err)
		}
		if exists {
			return errLinkAlreadyExists
		}
	}

	if _, err := context.Printf("Linking %q -> %q\n", to, from); err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}

	// make the parent folder
	if e := os.MkdirAll(parentTo, dirs.NewModBits); e != nil {
		return fmt.Errorf("%w: %w", errLinkUnknown, e)
	}

	// and make the symlink
	if e := os.Symlink(from, to); e != nil {
		return fmt.Errorf("%w: %w", errLinkUnknown, e)
	}

	return nil
}
