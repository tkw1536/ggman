package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Root is the 'ggman root' command.
//
// The 'ggman root' command prints the ggman root directory to standard output.
// It does not require the root directory to exist.
var Root program.Command = root{}

type root struct{}

func (root) BeforeRegister(program *program.Program) {}

func (root) Description() program.Description {
	return program.Description{
		Name:        "root",
		Description: "Print the ggman root folder. ",

		Environment: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (root) AfterParse() error {
	return nil
}

func (root) Run(context program.Context) error {
	context.Println(ggman.C2E(context).Root)
	return nil
}
