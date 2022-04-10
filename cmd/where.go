package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
)

// Where is the 'ggman where' command.
//
// When invoked, the ggman where command prints to standard output the location where the remote repository described by the first argument would be cloned to.
// This location is a subfolder of the directory outputted by 'ggman root'.
// Each segment of the path correesponding to a component of the repository url.
//
// This command does not perform any interactions with the remote repository or the local disk, in particular it does not require access to the remote repository or require it to be installed.
var Where ggman.Command = &where{}

type where struct {
	Positionals struct {
		URL string `required:"1-1" positional-arg-name:"URL" description:"Remote repository URL to use"`
	} `positional-args:"true"`
}

func (where) BeforeRegister(program *ggman.Program) {}

func (where) Description() ggman.Description {
	return ggman.Description{
		Command:     "where",
		Description: "Print the location where a repository would be cloned to",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (w *where) Run(context ggman.Context) error {
	localPath, err := context.Environment.Local(env.ParseURL(w.Positionals.URL))
	if err != nil {
		return err
	}
	context.Println(localPath)
	return nil
}
