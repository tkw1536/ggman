package cmd

import (
	"github.com/tkw1536/ggman/program"
)

// Comps is the 'ggman comps' command.
//
// When invoked, it prints the components of the first argument passed to it.
// Each component is printed on a seperate line of standard output.
var Comps program.Command = comps{}

type comps struct{}

func (comps) BeforeRegister(program *program.Program) {}

func (comps) Description() program.Description {
	return program.Description{
		Name:        "comps",
		Description: "Print the components of a URL",

		PosArgsMin: 1,
		PosArgsMax: 1,

		PosArgName: "URL",

		PosArgDescription: "URL to use",
	}
}

func (comps) AfterParse() error {
	return nil
}

func (comps) Run(context program.Context) error {
	for _, comp := range context.URLV(0).Components() {
		context.Println(comp)
	}

	return nil
}
