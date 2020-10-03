package cmd

import (
	"github.com/spf13/pflag"

	"github.com/tkw1536/ggman/program"
)

// Comps is the 'ggman comps' command.
//
// When invoked, it prints the components of the first argument passed to it.
// Each component is printed on a seperate line of standard output.
var Comps program.Command = comps{}

type comps struct{}

func (comps) Name() string {
	return "comps"
}

func (comps) Options(flagset *pflag.FlagSet) program.Options {
	return program.Options{
		MinArgs: 1,
		MaxArgs: 1,

		Metavar: "URI",

		UsageDescription: "Repository URI to print components of. ",
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
