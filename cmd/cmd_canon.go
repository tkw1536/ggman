package cmd

import (
	flag "github.com/spf13/pflag"

	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Canon is the 'ggman canon' command
var Canon program.Command = canon{}

type canon struct{}

func (canon) Name() string {
	return "canon"
}

func (canon) Options(flagset *flag.FlagSet) program.Options {
	return program.Options{
		MinArgs: 1,
		MaxArgs: 2,

		UsageDescription: "The URL of which to get the canonical location and an optional CANSPEC. ",
	}
}

func (canon) AfterParse() error {
	return nil
}

func (canon) Run(context program.Context) error {
	var file env.CanFile

	switch len(context.Argv) {
	case 1: // read the default CanFile
		if err := (&(context.Env)).LoadDefaultCANFILE(); err != nil { //TODO: This breaks test isolation

			return err
		}
	case 2: // use a custom CanLine
		file = []env.CanLine{{Pattern: "", Canonical: context.Argv[1]}}
	}

	// print out the canonical version of the file
	canonical := env.ParseURL(context.Argv[0]).CanonicalWith(file)
	context.Println(canonical)

	return nil
}
