package cmd

//spellchecker:words github cobra ggman
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
)

//spellchecker:words wrapcheck

func NewLicenseCommand() *cobra.Command {
	impl := new(license)

	cmd := &cobra.Command{
		Use:   "license",
		Short: "print license information about ggman and exit",
		Long:  "The license command prints to standard output legal notices about the ggman program.",
		Args:  cobra.NoArgs,

		RunE: impl.Exec,
	}

	return cmd
}

type license struct{}

func (license) Exec(cmd *cobra.Command, args []string) error {
	_, err := fmt.Fprintf(cmd.OutOrStdout(), stringLicenseInfo, ggman.License, ggman.Notices)
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	return nil
}

const stringLicenseInfo = `
ggman -- A golang script that can manage multiple git repositories locally
https://go.tkw01536.de/ggman

================================================================================
ggman is licensed under the terms of the MIT License:

%s
================================================================================

Furthermore, this executable may include code from the following projects:
%s
`
