//spellchecker:words ggman
package ggman

//spellchecker:words strconv time

//spellchecker:words ggman

// DefaultBuildVersion is the build version used when no information is available.
const DefaultBuildVersion = "v0.0.0-unknown"

// these private constants are set by the Makefile at build time.
var buildVersion string = DefaultBuildVersion

// BuildVersion is the current version of this program.
// This is only set by the ggman build process, and can not be found in this documentation.
//
// When the build version is not known, it is set to DefaultBuildVersion.
var BuildVersion string = buildVersion
