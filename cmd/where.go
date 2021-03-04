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

func (where) Name() string {
	return "where"
}

func (where) Options() program.Options {
	return program.Options{
		Description: "Print the location where a repository would be cloned to",

		MinArgs: 1,
		MaxArgs: 1,

		Metavar: "URL",

		UsageDescription: "Remote repository URL to use",

		Environment: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (where) AfterParse() error {
	return nil
}

func (where) Run(context program.Context) error {
	localPath := context.Local(context.URLV(0))
	context.Println(localPath)
	return nil
}
