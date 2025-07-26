package cmd

import (
	"errors"
	"fmt"

	"al.essio.dev/pkg/shellescape"
	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/ggman/git"
	"go.tkw01536.de/goprogram/exit"
)

//spellchecker:words nolint wrapcheck

var (
	errCloneInvalidDestFlags = exit.NewErrorWithCode(`invalid destination: "--to" and "--plain" may not be used together`, exit.ExitCommandArguments)
	errCloneInvalidDest      = exit.NewErrorWithCode("unable to determine local destination", exit.ExitGeneralArguments)
	errCloneLocalURI         = exit.NewErrorWithCode("invalid remote URI: invalid scheme, not a remote path", exit.ExitCommandArguments)
	errCloneAlreadyExists    = exit.NewErrorWithCode("unable to clone repository: another git repository already exists in target location", exit.ExitGeneric)
	errCloneNoArguments      = exit.NewErrorWithCode("external `git` not found, can not pass any additional arguments to `git clone`", exit.ExitGeneric)
	errCloneOther            = exit.NewErrorWithCode("", exit.ExitGeneric)
	errCloneNoComps          = errors.New("unable to find components of URI")
)

func NewCloneCommand() *cobra.Command {
	// Introduce variables for flags
	var (
		force bool
		local bool
		exact bool
		plain bool
		to    string
	)

	clone := &cobra.Command{
		Use:   "clone",                                                       // Set old command name
		Short: "clone a repository into a path described by \"ggman where\"", // Copy old description
		Args:  cobra.MinimumNArgs(1),                                         // Set appropriate argument validation

		PreRunE: func(cmd *cobra.Command, args []string) error {
			ggman.SetRequirements(cmd, &env.Requirement{
				NeedsRoot:    true,
				NeedsCanFile: true,
			})

			if local {
				plain = true
			}
			if plain && to != "" {
				return errCloneInvalidDestFlags
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			environment, err := ggman.GetEnv(cmd)
			if err != nil {
				return fmt.Errorf("failed to get environment: %w", err)
			}

			url := env.ParseURL(args[0])
			if url.IsLocal() {
				return fmt.Errorf("%q: %w", args[0], errCloneLocalURI)
			}

			// determine paths
			remote := args[0]
			if !exact {
				remote = environment.Canonical(url)
			}
			local, err := findDest(&environment, url, plain, to)
			if err != nil {
				return fmt.Errorf("%w: %w", errCloneInvalidDest, err)
			}

			// do the clone!
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Cloning %q into %q ...\n", remote, local); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
			switch err := environment.Git.Clone(streamFromCommand(cmd), remote, local, args[1:]...); {
			case err == nil:
				return nil
			case errors.Is(err, git.ErrCloneAlreadyExists):
				if force {
					if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Clone already exists in target location, done."); err != nil {
						return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
					}
					return nil
				}
				return errCloneAlreadyExists
			case errors.Is(err, git.ErrArgumentsUnsupported):
				return fmt.Errorf("%w: %v", errCloneNoArguments, shellescape.QuoteCommand(args[1:]))
			default:
				return fmt.Errorf("%w%w", errCloneOther, err)
			}
		},
	}

	// Add flags
	flags := clone.Flags()
	flags.BoolVarP(&force, "force", "f", false, "do not complain when a repository already exists in the target directory")
	flags.BoolVarP(&local, "local", "l", false, "alias of \"--plain\"")
	flags.BoolVarP(&exact, "exact-url", "e", false, "don't canonicalize URL before cloning and use exactly the passed URL")
	flags.BoolVar(&plain, "plain", false, "clone like a standard git would: into an appropriately named subdirectory of the current directory")
	flags.StringVarP(&to, "to", "t", "", "clone repository into specified directory")

	return clone
}

//nolint:wrapcheck // wrapped in the caller
func findDest(environment *env.Env, url env.URL, plain bool, to string) (string, error) {
	switch {
	case plain:
		comps := url.Components()
		if len(comps) == 0 {
			return "", errCloneNoComps
		}
		return environment.Abs(comps[len(comps)-1])
	case to != "":
		return environment.Abs(to)
	default:
		return environment.Local(url)
	}
}
