package ggman

import (
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram"
)

// URLV returns the ith parameters as a url.
func URLV[E any, P any, F any, R goprogram.Requirement[F]](context goprogram.Context[E, P, F, R], i int) env.URL {
	return env.ParseURL(context.Args.Pos[i])
}
