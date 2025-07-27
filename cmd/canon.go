package cmd

//spellchecker:words github cobra ggman pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words github cobra ggman goprogram exit

func NewCanonCommand() *cobra.Command {
	impl := new(canon)

	cmd := &cobra.Command{
		Use:   "canon URL [CANSPEC]",
		Short: "print the canonical version of a URL",
		Long: `The 'ggman canon' command prints to standard output the canonical version of the URL passed as the first argument.
		An optional second argument determines the CANSPEC to use for canonizing the URL.`,
		Args: cobra.RangeArgs(1, 2),

		PreRunE: PreRunE(impl),
		RunE:    impl.Exec,
	}

	return cmd
}

type canon struct {
	Positional struct {
		URL     env.URL
		CANSPEC string
	}
}

var (
	errCanonUnableCanFile = exit.NewErrorWithCode("unable to load default CANFILE", exit.ExitContext)
)

func (c *canon) AfterParse(cmd *cobra.Command, args []string) error {
	c.Positional.URL = env.ParseURL(args[0])
	if len(args) == 2 {
		c.Positional.CANSPEC = args[1]
	}
	return nil
}

func (c *canon) Exec(cmd *cobra.Command, args []string) error {
	var file env.CanFile

	if c.Positional.CANSPEC == "" {
		env, err := ggman.GetEnv(cmd, env.Requirement{})
		if err != nil {
			return fmt.Errorf("failed to get environment: %w", err)
		}

		if file, err = env.LoadDefaultCANFILE(); err != nil {
			return fmt.Errorf("%w: %w", errCanonUnableCanFile, err)
		}
	} else {
		file = []env.CanLine{{Pattern: "", Canonical: c.Positional.CANSPEC}}
	}

	// print out the canonical version of the file
	canonical := c.Positional.URL.CanonicalWith(file)

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), canonical); err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}
	return nil
}
