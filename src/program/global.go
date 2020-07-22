package program

import (
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/tkw1536/ggman/src/repos"
)

// GetRootOrPanic gets the default root folder or panics()
func GetRootOrPanic() (value string, err error) {
	value = os.Getenv("GGROOT")
	if len(value) == 0 {
		value, err = homedir.Expand("~/Projects")
	}

	return
}

// GetCanonOrPanic returns the default canon file or panics
func GetCanonOrPanic() (file repos.CanFile, err error) {
	return file, (&file).ReadDefault()
}
