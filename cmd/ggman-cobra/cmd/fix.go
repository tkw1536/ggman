package cmd

//spellchecker:words sync github cobra ggman goprogram exit
import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
)

func NewFixCommand() *cobra.Command {
	impl := new(fix)

	cmd := &cobra.Command{
		Use:   "fix",
		Short: "canonicalizes remote URLs for cloned repositories",
		Long:  `The 'ggman fix' command canonicalizes the urls of all remotes of a repository.`,
		Args:  cobra.NoArgs,

		PreRunE: PreRunE(impl),
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.Simulate, "simulate", "s", false, "do not perform any canonicalization, instead only print what would be done")

	return cmd
}

//spellchecker:words canonicalizes canonicalization nolint wrapcheck

type fix struct {
	Simulate bool
}

func (*fix) AfterParse(cmd *cobra.Command, args []string) error {
	return nil
}

func (fix) Description() ggman.Description {
	return ggman.Description{
		Command:     "fix",
		Description: "canonicalizes remote URLs for cloned repositories",
		Requirements: env.Requirement{
			NeedsRoot:    true,
			NeedsCanFile: true,
			AllowsFilter: true,
		},
	}
}

var errFixCustom = exit.NewErrorWithCode("", exit.ExitGeneric)

func (f *fix) Exec(cmd *cobra.Command, args []string) error {
	environment, err := ggman.GetEnv(cmd)
	if err != nil {
		return err
	}

	simulate := f.Simulate

	hasError := false

	var innerError error
	for _, repo := range environment.Repos(true) {
		var initialMessage sync.Once // send an initial log message to the user, once

		if e := environment.Git.UpdateRemotes(repo, func(url, remoteName string) (string, error) {
			canon := environment.Canonical(env.ParseURL(url))

			if url == canon {
				return url, nil
			}

			initialMessage.Do(func() {
				if !simulate {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Fixing remote of %q", repo); err != nil {
						innerError = fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
					}
				} else {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Simulate fixing remote of %q", repo); err != nil {
						innerError = fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
					}
				}
			})

			if innerError != nil {
				return "", innerError
			}

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Updating %s: %s -> %s\n", remoteName, url, canon); err != nil {
				return "", fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
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
