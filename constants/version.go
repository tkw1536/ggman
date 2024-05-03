//spellchecker:words constants
package constants

//spellchecker:words strconv time
import (
	"strconv"
	"time"
)

//spellchecker:words ggman

// DefaultBuildVersion is the build version used when no information is available.
const DefaultBuildVersion = "v0.0.0-unknown"

// these private constants are set by the Makefile at build time.
var buildTime string = "0"
var buildVersion string = DefaultBuildVersion

// BuildTime is the time this program was built.
// This is only set by the ggman build process, and can not be found in this documentation.
//
// When the build time is not known, it is set to 1970-01-01.
var BuildTime time.Time

func init() {
	// setup time properly
	buildTimeInt, err := strconv.ParseInt(buildTime, 0, 64)
	if err != nil {
		panic("constants.buildTime invalid")
	}
	BuildTime = time.Unix(buildTimeInt, 0).UTC()
}

// BuildVersion is the current version of this program.
// This is only set by the ggman build process, and can not be found in this documentation.
//
// When the build version is not known, it is set to DefaultBuildVersion.
var BuildVersion string = buildVersion
