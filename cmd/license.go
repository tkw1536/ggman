package cmd

//spellchecker:words ggman constants legal
import (
	"fmt"

	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/constants/legal"
)

//spellchecker:words nolint wrapcheck

// License is the 'ggman license' command.
//
// The license command prints to standard output legal notices about the ggman program.
var License ggman.Command = license{}

type license struct{}

func (license) Description() ggman.Description {
	return ggman.Description{
		Command:     "license",
		Description: "print license information about ggman and exit",
	}
}

func (license) Run(context ggman.Context) error {
	_, err := context.Printf(stringLicenseInfo, ggman.License, legal.Notices)
	if err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
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
