package constants

import (
	"strconv"
	"time"
)

// the following constants are related to versioning
// and are set at buildtime by the Makefile.
var buildTime string = "0"
var buildVersion string = "v0.0.0-unknown"

// BuildTime is the time this program was built
// When the build time is not known, it is set to 1970-01-01.
var BuildTime time.Time

func init() {
	buildTimeInt, err := strconv.ParseInt(buildTime, 0, 64)
	if err != nil {
		panic("constants.buildTime invalid")
	}
	BuildTime = time.Unix(buildTimeInt, 0).UTC()
}

// BuildVersion is the current version of this program
var BuildVersion string = buildVersion
