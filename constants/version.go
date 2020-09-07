package constants

import (
	"strconv"
	"time"
)

// these private constants are set by the Makefile at build time.
var buildTime string = "0"
var buildVersion string = "v0.0.0-unknown"

// BuildTime is the time this program was built.
// This is only set by the ggman build process, and can not be found in this documentation.
//
// When the build time is not known, it is set to 1970-01-01.
var BuildTime time.Time

func init() {
	buildTimeInt, err := strconv.ParseInt(buildTime, 0, 64)
	if err != nil {
		panic("constants.buildTime invalid")
	}
	BuildTime = time.Unix(buildTimeInt, 0).UTC()
}

// BuildVersion is the current version of this program.
// This is only set by the ggman build process, and can not be found in this documentation.
//
// When the build version is not known, it is set to "v0.0.0-unknown".
var BuildVersion string = buildVersion
