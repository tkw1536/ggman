package cmd

//spellchecker:words github ggman
import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
)

// Comps is the 'ggman comps' command.
//
// When invoked, it prints the components of the first argument passed to it.
// Each component is printed on a separate line of standard output.
var Comps ggman.Command = comps{}

type comps struct {
	Positional struct {
		URL env.URL `required:"1-1" positional-arg-name:"URL" description:"URL to use"`
	} `positional-args:"true"`
}

func (comps) Description() ggman.Description {
	return ggman.Description{
		Command:     "comps",
		Description: "print the components of a URL",
	}
}

func (c comps) Run(context ggman.Context) error {
	for _, comp := range c.Positional.URL.Components() {
		if _, err := context.Println(comp); err != nil {
			return ggman.ErrGenericOutput.WrapError(err)
		}
	}

	return nil
}
