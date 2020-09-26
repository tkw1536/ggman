package cmd

import (
	"github.com/spf13/pflag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Lsr is the 'ggman lsr' command
var Lsr program.Command = &lsr{}

type lsr struct {
	Canonical bool
}

func (lsr) Name() string {
	return "lsr"
}

func (l *lsr) Options(flagset *pflag.FlagSet) program.Options {
	flagset.BoolVarP(&l.Canonical, "canonical", "c", l.Canonical, "If provided, print the canonical URLs to repositories. ")
	return program.Options{
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

func (l lsr) Run(context program.Context) error {
	var lines env.CanFile
	if l.Canonical {
		if err := (&context.Env).LoadDefaultCANFILE(); err != nil {
			return errInvalidCanfile
		}
		lines = context.CanFile
	}

	// and print them
	for _, repo := range context.Repos() {
		remote, err := context.Git.GetRemote(repo)
		if err != nil {
			continue
		}
		if l.Canonical {
			remote = env.ParseURL(remote).CanonicalWith(lines)
		}
		context.Println(remote)
	}

	return nil
}
