// Package ggman serves as the implementation of the ggman program.
// See documentation of the ggman command as an entry point into the documentation.
//
// Note that this package and it's subpackages are not intended to be consumed by other go packages.
// The public interface of the ggman is defined only by the ggman executable.
// This package is not considered part of the public interface as such and not subject to Semantic Versioning.
//
// The top-level ggman package is considered to be stand-alone, and (with the exception of 'util') does not directly depend on any of its' subpackages.
// As such it can be safely used by any subpackage without cyclic imports.
package ggman

import (
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram"
)

type ggmanEnv = env.Env
type ggmanParameters = env.Parameters
type ggmanRequirements = env.Requirement
type ggmanFlags = env.Flags

// Program is the type of the ggman Program
type Program = goprogram.Program[ggmanEnv, ggmanParameters, ggmanFlags, ggmanRequirements]

// Command is the type of a ggman Command
type Command = goprogram.Command[ggmanEnv, ggmanParameters, ggmanFlags, ggmanRequirements]

// Context is type type of a Context passed to ggman command
type Context = goprogram.Context[ggmanEnv, ggmanParameters, ggmanFlags, ggmanRequirements]

// Arguments is the type of ggman Arguments
type Arguments = goprogram.Arguments[ggmanFlags]

// Description is the type of descriptuions of a ggman command
type Description = goprogram.Description[ggmanFlags, ggmanRequirements]
