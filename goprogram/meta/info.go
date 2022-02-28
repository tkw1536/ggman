package meta

import (
	"fmt"
	"runtime"
	"time"
)

// Info holds information about a program
type Info struct {
	BuildVersion string
	BuildTime    time.Time

	Executable  string // Name of the main executable of the program
	Description string // Description of the program
}

// FmtVersion formats version information about the current version
// It returns a string that should be presented to users.
func (info Info) FmtVersion() string {
	return fmt.Sprintf("%s version %s, built %s, using %s", info.Executable, info.BuildVersion, info.BuildTime, runtime.Version())
}
