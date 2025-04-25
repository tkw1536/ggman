package cmd

//spellchecker:words github ggman goprogram exit
import (
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
)

// this files contains errors shared by multiple commands.

var errUnableLocalPath = exit.Error{
	ExitCode: env.ExitInvalidRepo,
	Message:  "failed to get local path",
}
