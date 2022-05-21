package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/goprogram/parser"
)

// Complete is the 'ggman _complete' command.
// It provides tab completion
var Complete ggman.Command = complete{}

type complete struct {
	Positional struct {
		Rest []string `positional-arg-name:"ARG" description:"Arguments to tab complete"`
	} `positional-args:"true"`
}

func (complete) Description() ggman.Description {
	return ggman.Description{
		Command:     "_complete",
		Description: "tab complete a partial command line",
		ParserConfig: parser.Config{
			IncludeUnknown: true,
		},
	}
}

func (c complete) Run(context ggman.Context) error {
	res, err := context.Program.Complete(c.Positional.Rest)
	if err != nil {
		return err
	}

	for _, r := range res {
		context.Println(r)
	}
	return nil
}
