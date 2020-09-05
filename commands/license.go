package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/program"
)

// LicenseCommand is the entry point for the license command
func LicenseCommand(runtime *program.SubRuntime) (retval int, err string) {
	fmt.Printf(constants.StringLicenseInfo, constants.StringLicenseText, constants.StringLicenseNotices)

	return
}
