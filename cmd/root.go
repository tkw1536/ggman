package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
)

// Root is the 'ggman root' command.
//
// The 'ggman root' command prints the ggman root directory to standard output.
// It does not require the root directory to exist.
var Root ggman.Command = root{}

type root struct{}

func (root) BeforeRegister(program *ggman.Program) {}

func (root) Description() ggman.Description {
	return ggman.Description{
		Command:     "root",
		Description: "Print the ggman root folder. ",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (root) AfterParse() error {
	return nil
}

func (root) Run(context ggman.Context) error {
	context.Println(context.Environment.Root)
	return nil
}
