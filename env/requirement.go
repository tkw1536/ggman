package env

import (
	"reflect"

	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/program/exit"
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

var errTakesNoArgument = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes no '%s' argument. ",
}

func (req Requirement) Validate(args program.Arguments) error {
	if req.AllowsFilter { // no checking needed!
		return nil
	}

	// check the value of the arguments struct
	aVal := reflect.ValueOf(args.Flags)

	for _, fIndex := range flagsIndexes {
		v := aVal.FieldByIndex(fIndex)

		if !v.IsZero() { // flag was set iff it is non-zero
			tp := flagsType.FieldByIndex(fIndex) // needed for the error message only!
			return errTakesNoArgument.WithMessageF(args.Command, "--"+tp.Tag.Get("long"))
		}
	}

	return nil
}

// reflect access to the arguments type
var flagsType reflect.Type = reflect.TypeOf((*program.Flags)(nil)).Elem()
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
