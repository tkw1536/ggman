package cmd

//spellchecker:words sync github cobra ggman internal pkglib exit
import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words canonicalize

func NewFixCommand() *cobra.Command {
	impl := new(fix)

	cmd := &cobra.Command{
		Use:   "fix",
		Short: "Canonicalize remote URLs for cloned repositories",
		Long: `The 'ggman fix' command canonicalizes the urls of all remotes of a repository.

This updates remotes of all matching repositories to their canonical form using the CANFILE.
Optionally, you can pass a '--simulate' argument to 'ggman fix'.
Instead of storing any urls, it will only print what is being done to STDOUT.`,
		Args: cobra.NoArgs,

		RunE: impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.Simulate, "simulate", "s", false, "do not perform any canonicalization, instead only print what would be done")
	flags.BoolVarP(&impl.PruneRemotes, "prune-remotes", "p", false, "prune unused remotes")

	return cmd
}

//spellchecker:words canonicalizes canonicalization wrapcheck

type fix struct {
	Simulate     bool
	PruneRemotes bool
}

var errFixCustom = exit.NewErrorWithCode("", env.ExitGeneric)

func (f *fix) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		NeedsRoot:    true,
		NeedsCanFile: true,
		AllowsFilter: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	simulate := f.Simulate

	hasError := false

	var innerError error
	for _, repo := range environment.Repos(cmd.Context(), true) {
		var initialMessage sync.Once // send an initial log message to the user, once

		// first prune unused remotes if requested (so we don't need to update them)
		if f.PruneRemotes {
			remotes, e := environment.Git.FindUnusedRemotes(cmd.Context(), repo)
			if e != nil {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), e.Error()) // no way to report error
				hasError = true
			}

			// iterate over the remotes
			for remote := range remotes {
				if !simulate {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Removing unused remote %q from %q\n", remote, repo); err != nil {
						_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err.Error())
						hasError = true
						continue
					}
					if e := environment.Git.DeleteRemotes(cmd.Context(), repo, remote); e != nil {
						_, _ = fmt.Fprintln(cmd.ErrOrStderr(), e.Error()) // no way to report error
						hasError = true
					}
				} else {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Found unused remote %q in %q\n", remote, repo); err != nil {
						_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err.Error()) // no way to report error
						hasError = true
					}
				}
			}
		}

		if e := environment.Git.UpdateRemotes(cmd.Context(), repo, func(url, remoteName string) (string, error) {
			canon := environment.Canonical(env.ParseURL(url))

			if url == canon {
				return url, nil
			}

			initialMessage.Do(func() {
				if !simulate {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Fixing remote of %q\n", repo); err != nil {
						innerError = fmt.Errorf("%w: %w", errGenericOutput, err)
					}
				} else {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Simulate fixing remote of %q\n", repo); err != nil {
						innerError = fmt.Errorf("%w: %w", errGenericOutput, err)
					}
				}
			})

			if innerError != nil {
				return "", innerError
			}

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Updating %s: %s -> %s\n", remoteName, url, canon); err != nil {
				return "", fmt.Errorf("%w: %w", errGenericOutput, err)
			}

			// either return the canonical url, or (if we're simulating) the old url
			if simulate {
				return url, nil
			}

			return canon, nil
		}); e != nil {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), e.Error()) // no way to report error
			hasError = true
		}
	}

	// if we had an error, indicate that to the user
	if hasError {
		return errFixCustom
	}

	// and finish
	return nil
}
