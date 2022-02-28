package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/goprogram/meta"
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
		Command:     "canon",
		Description: "Print the canonical version of a URL",

		Positional: meta.Positional{
			Description: "URL and optional CANSPEC of repository",
			Min:         1,
			Max:         2,
		},
	}
}

func (canon) AfterParse() error {
	return nil
}

func (canon) Run(context ggman.Context) error {
	var file env.CanFile

	switch len(context.Args.Pos) {
	case 1: // read the default CanFile
		var err error
		if file, err = context.Environment.LoadDefaultCANFILE(); err != nil {
			return err
		}
	case 2: // use a custom CanLine
		file = []env.CanLine{{Pattern: "", Canonical: context.Args.Pos[1]}}
	}

	// print out the canonical version of the file
	canonical := ggman.URLV(context, 0).CanonicalWith(file)
	context.Println(canonical)

	return nil
}
