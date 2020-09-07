package cmd

import (
	flag "github.com/spf13/pflag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Lsr is the 'ggman lsr' command
var Lsr program.Command = lsr{}

type lsr struct{}

func (lsr) Name() string {
	return "lsr"
}

func (lsr) Options(flagset *flag.FlagSet) program.Options {
	return program.Options{
		FlagValue:       "--canonical",
		FlagDescription: "If provided, print the canonical URLs to repositories. ",

		Environment: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

func (lsr) AfterParse() error {
	return nil
}

var errInvalidCanfile = ggman.Error{
	Message:  "Invalid CANFILE found. ",
	ExitCode: ggman.ExitInvalidEnvironment,
}

func (lsr) Run(context program.Context) error {
	shouldCanon := context.Flag

	var lines env.CanFile
	if shouldCanon {
		if err := (&context.Env).LoadDefaultCANFILE(); err != nil {
			return errInvalidCanfile
		}
	}

	// and print them
	for _, repo := range context.Repos() {
		remote, err := context.Git.GetRemote(repo)
		if err != nil {
			continue
		}
		if shouldCanon {
			remote = context.ParseURL(remote).CanonicalWith(lines)
		}
		context.Println(remote)
	}

	return nil
}
