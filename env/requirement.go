package env

import (
	"reflect"

	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/program/usagefmt"
)

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

// AllowsOption checks if the provided option is allowed by this option
func (req Requirement) AllowsOption(opt usagefmt.Opt) bool {
	return req.AllowsFilter
}

func (req Requirement) Validate(args program.Arguments[Flags]) error {
	return program.ValidateAllowedOptions[Flags](req, args)
}

// reflect access to the arguments type
var flagsType reflect.Type = reflect.TypeOf((*Flags)(nil)).Elem()
var flagsIndexes [][]int // indexes of all the options

func init() {
	// iterate over the fields of the type
	fieldCount := flagsType.NumField()
	for i := 0; i < fieldCount; i++ {
		field := flagsType.Field(i)

		// it's a long filter name
		flagsIndexes = append(flagsIndexes, field.Index)
	}
}
