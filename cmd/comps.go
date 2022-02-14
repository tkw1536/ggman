package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/program/meta"
)

// Comps is the 'ggman comps' command.
//
// When invoked, it prints the components of the first argument passed to it.
// Each component is printed on a seperate line of standard output.
var Comps ggman.Command = comps{}

type comps struct{}

func (comps) BeforeRegister(program *ggman.Program) {}

func (comps) Description() ggman.Description {
	return ggman.Description{
		Command:     "comps",
		Description: "Print the components of a URL",

		Positional: meta.Positional{
			Value:       "URL",
			Description: "URL to use",

			Min: 1,
			Max: 1,
		},
	}
}

func (comps) AfterParse() error {
	return nil
}

func (comps) Run(context ggman.Context) error {
	for _, comp := range ggman.URLV(context, 0).Components() {
		context.Println(comp)
	}

	return nil
}
