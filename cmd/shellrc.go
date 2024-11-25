package cmd

//spellchecker:words embed github ggman
import (
	_ "embed"

	"github.com/tkw1536/ggman"
)

//spellchecker:words shellrc

// Shellrc is the 'ggman shellrc' command.
//
// The 'ggman shellrc' command prints aliases to be used for shell profiles in conjunction with ggman.
var Shellrc ggman.Command = shellrc{}

type shellrc struct{}

func (shellrc) Description() ggman.Description {
	return ggman.Description{
		Command:     "shellrc",
		Description: "print additional aliases to be used in shell profiles in conjunction with ggman",
	}
}

//go:embed shellrc.sh
var shellrcSh string

func (shellrc) Run(context ggman.Context) error {
	_, err := context.Printf("%s", shellrcSh)
	return ggman.ErrGenericOutput.WrapError(err)
}
