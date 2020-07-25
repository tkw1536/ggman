package constants

import (
	"strconv"
	"time"
)

var (
	buildTime string

	//BuildTime is the time this program was built or the string 'unknown'
	BuildTime string

	// BuildVersion is the git version of this program or the srting 'unknown'
	BuildVersion string
)

func init() {
	if buildTime == "" {
		BuildTime = "unknown"
	} else {
		i, err := strconv.ParseInt(buildTime, 10, 64)
		if err != nil {
			panic(err)
		}
		BuildTime = time.Unix(i, 0).String()
	}
	if BuildVersion == "" {
		BuildVersion = "unknown"
	}
}
