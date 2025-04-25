package env

//spellchecker:words github goprogram meta
import (
	"github.com/tkw1536/goprogram"
	"github.com/tkw1536/goprogram/meta"
)

//spellchecker:words nolint wrapcheck

// Requirement represents a set of requirements on the Environment.
type Requirement struct {
	// Does the environment require a root directory?
	NeedsRoot bool

	// Does the environment allow filtering?
	// AllowsFilter implies NeedsRoot.
	AllowsFilter bool

	// Does the environment require a CanFile?
	NeedsCanFile bool
}

// AllowsFlag checks if the provided option is allowed by this option.
func (req Requirement) AllowsFlag(flag meta.Flag) bool {
	return req.AllowsFilter
}

func (req Requirement) Validate(args goprogram.Arguments[Flags]) error {
	return goprogram.ValidateAllowedFlags(req, args) // nolint:wrapcheck
}
