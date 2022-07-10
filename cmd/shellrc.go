package cmd

import (
	_ "embed"

	"github.com/tkw1536/ggman"
)

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
	context.Printf("%s", shellrcSh)
	return nil
}
