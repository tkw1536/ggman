package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
)

// Canon is the 'ggman canon' command.
//
// The 'ggman canon' command prints to standard output the canonical version of the URL passed as the first argument.
// An optional second argument determines the CANSPEC to use for canonizing the URL.
var Canon ggman.Command = &canon{}

type canon struct{}

func (canon) BeforeRegister(program *ggman.Program) {}

func (canon) Description() ggman.Description {
	return ggman.Description{
		Name:        "canon",
		Description: "Print the canonical version of a URL",

		PosArgsMin: 1,
		PosArgsMax: 2,

		PosArgDescription: "URL and optional CANSPEC of repository",
	}
}

func (canon) AfterParse() error {
	return nil
}

func (canon) Run(context ggman.Context) error {
	var file env.CanFile

	switch len(context.Args.Arguments.Pos) {
	case 1: // read the default CanFile
		var err error
		if file, err = context.Runtime().LoadDefaultCANFILE(); err != nil {
			return err
		}
	case 2: // use a custom CanLine
		file = []env.CanLine{{Pattern: "", Canonical: context.Args.Arguments.Pos[1]}}
	}

	// print out the canonical version of the file
	canonical := ggman.URLV(context, 0).CanonicalWith(file)
	context.Println(canonical)

	return nil
}
