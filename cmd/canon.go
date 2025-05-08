package cmd

//spellchecker:words github ggman goprogram exit
import (
	"fmt"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
)

//spellchecker:words CANSPEC CANFILE nolint wrapcheck

// Canon is the 'ggman canon' command.
//
// The 'ggman canon' command prints to standard output the canonical version of the URL passed as the first argument.
// An optional second argument determines the CANSPEC to use for canonizing the URL.
var Canon ggman.Command = canon{}

type canon struct {
	Positional struct {
		URL     env.URL `description:"URL of the repository"     positional-arg-name:"URL"     required:"1-1"`
		CANSPEC string  `description:"CANSPEC of the repository" positional-arg-name:"CANSPEC"`
	} `positional-args:"true"`
}

func (canon) Description() ggman.Description {
	return ggman.Description{
		Command:     "canon",
		Description: "print the canonical version of a URL",
	}
}

var (
	errCanonUnableCanFile = exit.NewErrorWithCode("unable to load default CANFILE", exit.ExitContext)
)

func (c canon) Run(context ggman.Context) error {
	var file env.CanFile

	if c.Positional.CANSPEC == "" {
		var err error
		if file, err = context.Environment.LoadDefaultCANFILE(); err != nil {
			return fmt.Errorf("%w: %w", errCanonUnableCanFile, err)
		}
	} else {
		file = []env.CanLine{{Pattern: "", Canonical: c.Positional.CANSPEC}}
	}

	// print out the canonical version of the file
	canonical := c.Positional.URL.CanonicalWith(file)
	_, err := context.Println(canonical)
	if err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}
	return nil
}
