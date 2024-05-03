package cmd

//spellchecker:words github ggman
import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
)

//spellchecker:words CANSPEC CANFILE

// Canon is the 'ggman canon' command.
//
// The 'ggman canon' command prints to standard output the canonical version of the URL passed as the first argument.
// An optional second argument determines the CANSPEC to use for canonizing the URL.
var Canon ggman.Command = canon{}

type canon struct {
	Positional struct {
		URL     env.URL `required:"1-1" positional-arg-name:"URL" description:"URL of the repository"`
		CANSPEC string  `positional-arg-name:"CANSPEC" description:"CANSPEC of the repository"`
	} `positional-args:"true"`
}

func (canon) Description() ggman.Description {
	return ggman.Description{
		Command:     "canon",
		Description: "print the canonical version of a URL",
	}
}

func (c canon) Run(context ggman.Context) error {
	var file env.CanFile

	if c.Positional.CANSPEC == "" {
		var err error
		if file, err = context.Environment.LoadDefaultCANFILE(); err != nil {
			return err
		}
	} else {
		file = []env.CanLine{{Pattern: "", Canonical: c.Positional.CANSPEC}}
	}

	// print out the canonical version of the file
	canonical := c.Positional.URL.CanonicalWith(file)
	context.Println(canonical)

	return nil
}
