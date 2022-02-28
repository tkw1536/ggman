package ggman

import (
	"github.com/tkw1536/ggman/env"
	program "github.com/tkw1536/ggman/goprogram"
)

// URLV returns the ith parameters as a url.
func URLV[R any, P any, Flags any, Requirements program.Requirement[Flags]](context program.Context[R, P, Flags, Requirements], i int) env.URL {
	return env.ParseURL(context.Args.Pos[i])
}
