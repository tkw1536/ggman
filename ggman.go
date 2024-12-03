// Package ggman serves as the implementation of the ggman program.
// See documentation of the ggman command as an entry point into the documentation.
//
// Note that this package and it's sub-packages are not intended to be consumed by other go packages.
// The public interface of the ggman is defined only by the ggman executable.
// This package is not considered part of the public interface as such and not subject to Semantic Versioning.
//
// The top-level ggman package is considered to be stand-alone, and (with the exception of 'util') does not directly depend on any of its' sub-packages.
// As such it can be safely used by any sub-package without cyclic imports.
//
//spellchecker:words ggman
package ggman

//spellchecker:words github ggman goprogram exit
import (
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram"
	"github.com/tkw1536/goprogram/exit"
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

// Description is the type of descriptions of a ggman command
type Description = goprogram.Description[ggmanFlags, ggmanRequirements]

// ErrGenericOutput indicates that a generic output error occurred
var ErrGenericOutput = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "Unknown Output Error",
}
