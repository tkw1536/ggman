package cmd

import (
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Where is the 'ggman where' command.
//
// When invoked, the ggman where command prints to standard output the location where the remote repository described by the first argument would be cloned to.
// This location is a subfolder of the directory outputted by 'ggman root'.
// Each segment of the path correesponding to a component of the repository url.
//
// This command does not perform any interactions with the remote repository or the local disk, in particular it does not require access to the remote repository or require it to be installed.
var Where program.Command = where{}

type where struct{}

func (where) BeforeRegister(program *program.Program) {}

func (where) Description() program.Description {
	return program.Description{
		Name:        "where",
		Description: "Print the location where a repository would be cloned to",

		PosArgsMin: 1,
		PosArgsMax: 1,

		PosArgName: "URL",

		PosArgDescription: "Remote repository URL to use",

		Environment: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (where) AfterParse() error {
	return nil
}

func (where) Run(context program.Context) error {
	localPath, err := context.Env.Local(context.URLV(0))
	if err != nil {
		return err
	}
	context.Println(localPath)
	return nil
}
