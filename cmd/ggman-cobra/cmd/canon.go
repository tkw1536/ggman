package cmd

//spellchecker:words github cobra ggman goprogram exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
)

//spellchecker:words CANSPEC CANFILE nolint wrapcheck

var (
	errCanonUnableCanFile = exit.NewErrorWithCode("unable to load default CANFILE", exit.ExitContext)
)

func NewCanonCommand() *cobra.Command {
	var (
		URL     env.URL
		CANSPEC string
	)

	canon := &cobra.Command{
		Use:   "canon URL [CANSPEC]",
		Short: "print the canonical version of a URL",

		Args: cobra.RangeArgs(1, 2),

		PreRunE: func(cmd *cobra.Command, args []string) error {
			URL = env.ParseURL(args[0])
			if len(args) >= 2 {
				CANSPEC = args[1]
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			var file env.CanFile

			if CANSPEC == "" {
				env, err := ggman.GetEnv(cmd)
				if err != nil {
					return fmt.Errorf("failed to get environment: %w", err)
				}

				if file, err = env.LoadDefaultCANFILE(); err != nil {
					return fmt.Errorf("%w: %w", errCanonUnableCanFile, err)
				}
			} else {
				file = []env.CanLine{{Pattern: "", Canonical: CANSPEC}}
			}

			// print out the canonical version of the file
			canonical := URL.CanonicalWith(file)
			_, err := fmt.Fprintln(cmd.OutOrStdout(), canonical)
			if err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
			return nil
		},
	}

	return canon
}
