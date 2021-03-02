package cmd

import (
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Canon is the 'ggman canon' command.
//
// The 'ggman canon' command prints to standard output the canonical version of the URL passed as the first argument.
// An optional second argument determines the CANSPEC to use for canonizing the URL.
var Canon program.Command = &canon{}

type canon struct{}

func (canon) Name() string {
	return "canon"
}

func (canon) Options() program.Options {
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

	switch len(context.Args) {
	case 1: // read the default CanFile
		if err := (&(context.Env)).LoadDefaultCANFILE(); err != nil {
			return err
		}
		file = context.Env.CanFile
	case 2: // use a custom CanLine
		file = []env.CanLine{{Pattern: "", Canonical: context.Args[1]}}
	}

	// print out the canonical version of the file
	canonical := env.ParseURL(context.Args[0]).CanonicalWith(file)
	context.Println(canonical)

	return nil
}
