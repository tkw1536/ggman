package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/constants/legal"
	"github.com/tkw1536/ggman/program"
)

// License is the 'ggman license' command.
//
// The license command prints to standard output legal notices about the ggman program.
var License program.Command = license{}

type license struct{}

func (license) BeforeRegister(program *program.Program) {}

func (license) Description() program.Description {
	return program.Description{
		Name:        "license",
		Description: "Print license information about ggman and exit. ",
	}
}

func (license) AfterParse() error {
	return nil
}

func (license) Run(context program.Context) error {
	context.Printf(stringLicenseInfo, ggman.License, legal.Notices)
	return nil
}

const stringLicenseInfo = `
ggman -- A golang script that can manage multiple git repositories locally
https://github.com/tkw1536/ggman

================================================================================
ggman is licensed under the terms of the MIT License:

%s
================================================================================

Furthermore, this executable may include code from the following projects:
%s
`
