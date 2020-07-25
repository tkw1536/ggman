package commands

import (
	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/repos"
)

// HereCommand prints the current git repository root under the current path
func HereCommand(runtime *program.SubRuntime) (retval int, err string) {
	root, treePath := repos.Here(".", runtime.Root)
	if root == "" {
		return constants.ErrorInvalidRepo, constants.StringOutsideRepository
	}

	println(root)

	if runtime.Flag {
		println(treePath)
	}

	return
}
