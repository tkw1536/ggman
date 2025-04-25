package cmd

//spellchecker:words path filepath github ggman goprogram exit pkglib
import (
	"fmt"
	"path/filepath"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/fsx"
)

//spellchecker:words positionals nolint wrapcheck

// FindFile is the 'ggman find-file' command.
//
// The 'find-file' command lists all repositories that currently contain a file or directory with the provided name.
// The provided path may be relative to the root of the repository.
//
//	--exit-code
//
// When provided, exit with code 1 if no repositories are found.
//
//	--print-file
//
// Instead of listing the repository paths, print the filepath instead.
var FindFile ggman.Command = findFile{}

type findFile struct {
	Positionals struct {
		Path string `description:"name (or path) file to find" positional-arg-name:"PATH" required:"1-1"`
	} `positional-args:"true"`
	PrintFilePath bool `description:"instead of printing the repository paths, print the file paths"        long:"print-file" short:"p"`
	ExitCode      bool `description:"exit with status code 1 when no repositories with provided file exist" long:"exit-code"  short:"e"`
}

func (findFile) Description() ggman.Description {
	return ggman.Description{
		Command:     "find-file",
		Description: "list repositories containing a specific file",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (f findFile) AfterParse() error {
	if !filepath.IsLocal(f.Positionals.Path) {
		return errFindFileNotLocal
	}
	return nil
}

var errFindFileCustom = exit.Error{
	ExitCode: exit.ExitGeneric,
}

var errFindFileNotLocal = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "path argument is not a local path",
}

func (f findFile) Run(context ggman.Context) error {
	foundRepo := false
	for _, repo := range context.Environment.Repos(true) {
		candidate := filepath.Join(repo, f.Positionals.Path)
		ok, err := fsx.Exists(candidate)
		if err != nil {
			panic(err)
		}
		if !ok {
			continue
		}

		foundRepo = true
		if f.PrintFilePath {
			if _, err := context.Println(candidate); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
		} else {
			if _, err := context.Println(repo); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
		}
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if f.ExitCode && !foundRepo {
		return errFindFileCustom
	}

	return nil
}
