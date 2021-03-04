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

import _ "embed" // to include the license

// License contains the terms the ggman program is licensed under.
//go:embed LICENSE
var License string
