package cmd

//spellchecker:words path filepath ggman goprogram exit pkglib
import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/pkglib/fsx"
)

func NewFindFileCommand() *cobra.Command {
	impl := new(findFile)

	cmd := &cobra.Command{
		Use:   "find-file PATH",
		Short: "list repositories containing a specific file",
		Long: `The 'find-file' command lists all repositories that currently contain a file or directory with the provided name.
The provided path may be relative to the root of the repository.`,
		Args: cobra.ExactArgs(1),

		PreRunE: PreRunE(impl),
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.PrintFilePath, "print-file", "p", false, "instead of printing the repository paths, print the file paths")
	flags.BoolVarP(&impl.ExitCode, "exit-code", "e", false, "exit with status code 1 when no repositories with provided file exist")

	return cmd
}

//spellchecker:words positionals nolint wrapcheck

type findFile struct {
	Positionals struct {
		Path string
	}
	PrintFilePath bool
	ExitCode      bool
}

func (*findFile) Description() ggman.Description {
	return ggman.Description{
		Command:     "find-file",
		Description: "list repositories containing a specific file",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (f *findFile) AfterParse(cmd *cobra.Command, args []string) error {
	f.Positionals.Path = args[0]
	if !filepath.IsLocal(f.Positionals.Path) {
		return errFindFileNotLocal
	}
	return nil
}

var (
	errFindFileCustom   = exit.NewErrorWithCode("", exit.ExitGeneric)
	errFindFileNotLocal = exit.NewErrorWithCode("path argument is not a local path", exit.ExitCommandArguments)
)

func (f *findFile) Exec(cmd *cobra.Command, args []string) error {
	environment, err := ggman.GetEnv(cmd)
	if err != nil {
		return err
	}
	foundRepo := false
	for _, repo := range environment.Repos(true) {
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
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), candidate); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
		} else {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), repo); err != nil {
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
