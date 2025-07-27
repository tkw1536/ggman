package cmd

//spellchecker:words ggman constants legal
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/constants/legal"
)

//spellchecker:words nolint wrapcheck

// NewLicenseCommand creates the 'ggman license' command.
//
// The license command prints to standard output legal notices about the ggman program.
func NewLicenseCommand() *cobra.Command {
	impl := new(license)

	cmd := &cobra.Command{
		Use:   "license",
		Short: "print license information about ggman and exit",
		Long:  "The license command prints to standard output legal notices about the ggman program.",
		Args:  cobra.NoArgs,

		PreRunE: PreRunE(impl),
		RunE:    impl.Exec,
	}

	return cmd
}

type license struct{}

func (license) Description() ggman.Description {
	return ggman.Description{
		Command:     "license",
		Description: "print license information about ggman and exit",
	}
}

func (*license) AfterParse(cmd *cobra.Command, args []string) error {
	return nil
}

func (license) Exec(cmd *cobra.Command, args []string) error {
	_, err := fmt.Fprintf(cmd.OutOrStdout(), StringLicenseInfo, ggman.License, legal.Notices)
	if err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}
	return nil
}

// TODO: Make this private again?
const StringLicenseInfo = `
ggman -- A golang script that can manage multiple git repositories locally
https://go.tkw01536.de/ggman

================================================================================
ggman is licensed under the terms of the MIT License:

%s
================================================================================

Furthermore, this executable may include code from the following projects:
%s
`
