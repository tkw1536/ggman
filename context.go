package ggman

import (
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// URLV returns the ith parameters as a url.
func URLV[Runtime any, Parameters any, Requirements any](context program.Context[Runtime, Parameters, Requirements], i int) env.URL {
	return env.ParseURL(context.Args[i])
}
