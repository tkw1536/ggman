package cmd

//spellchecker:words github ggman
import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
)

//spellchecker:words positionals

// Where is the 'ggman where' command.
//
// When invoked, the ggman where command prints to standard output the location where the remote repository described by the first argument would be cloned to.
// This location is a subfolder of the directory outputted by 'ggman root'.
// Each segment of the path corresponding to a component of the repository url.
//
// This command does not perform any interactions with the remote repository or the local disk, in particular it does not require access to the remote repository or require it to be installed.
var Where ggman.Command = where{}

type where struct {
	Positionals struct {
		URL string `description:"remote repository URL to use" positional-arg-name:"URL" required:"1-1"`
	} `positional-args:"true"`
}

func (where) Description() ggman.Description {
	return ggman.Description{
		Command:     "where",
		Description: "print the location where a repository would be cloned to",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (w where) Run(context ggman.Context) error {
	localPath, err := context.Environment.Local(env.ParseURL(w.Positionals.URL))
	if err != nil {
		return err
	}
	_, err = context.Println(localPath)
	return ggman.ErrGenericOutput.WrapError(err)
}
