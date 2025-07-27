// Package ggman serves as the implementation of the ggman program.
// See documentation of the ggman command as an entry point into the documentation.
//
// Note that this package and it's sub-packages are not intended to be consumed by other go packages.
// The public interface of the ggman is defined only by the ggman executable.
// This package is not considered part of the public interface as such and not subject to Semantic Versioning.
//
// The top-level ggman package is considered to be stand-alone, and (with the exception of 'env') does not directly depend on any of its' sub-packages.
// As such it can be safely used by any sub-package without cyclic imports.
//
//spellchecker:words ggman
package ggman

//spellchecker:words context github cobra ggman pkglib exit
