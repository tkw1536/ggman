package ggman

import (
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// URLV returns the ith parameters as a url.
func URLV[Runtime any, Parameters any, Flags any, Requirements program.Requirement[Flags]](context program.Context[Runtime, Parameters, Flags, Requirements], i int) env.URL {
	return env.ParseURL(context.Args.Arguments.Pos[i])
}
