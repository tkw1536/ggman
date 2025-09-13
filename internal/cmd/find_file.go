package cmd

//spellchecker:words path filepath github cobra ggman internal pkglib exit
import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/fsx"
)

func NewFindFileCommand() *cobra.Command {
	impl := new(findFile)

	cmd := &cobra.Command{
		Use:   "find-file PATH",
		Short: "List repositories containing a specific file",
		Long: `Find-file lists all repositories that currently contain a file or directory with the provided name.
The provided path may be relative to the root of the repository.`,
		Args: cobra.ExactArgs(1),

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.PrintFilePath, "print-file", "p", false, "instead of printing the repository paths, print the file paths")
	flags.BoolVarP(&impl.ExitCode, "exit-code", "e", false, "exit with status code 1 when no repositories with provided file exist")

	return cmd
}

//spellchecker:words positionals wrapcheck

type findFile struct {
	Positionals struct {
		Path string
	}
	PrintFilePath bool
	ExitCode      bool
}

func (f *findFile) ParseArgs(cmd *cobra.Command, args []string) error {
	f.Positionals.Path = args[0]
	if !filepath.IsLocal(f.Positionals.Path) {
		return errFindFileNotLocal
	}
	return nil
}

var (
	errFindFileCustom   = exit.NewErrorWithCode("", env.ExitGeneric)
	errFindFileNotLocal = exit.NewErrorWithCode("path argument is not a local path", env.ExitCommandArguments)
)

func (f *findFile) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{NeedsRoot: true})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}
	foundRepo := false
	for _, repo := range environment.Repos(cmd.Context(), true) {
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
				return fmt.Errorf("%w: %w", errGenericOutput, err)
			}
		} else {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), repo); err != nil {
				return fmt.Errorf("%w: %w", errGenericOutput, err)
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
