package env

import (
	"reflect"

	"github.com/tkw1536/ggman/internal/text"
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

var errTakesNoArgument = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes no '%s' argument. ",
}

func (req Requirement) Validate(args program.Arguments) error {
	if req.AllowsFilter { // no checking needed!
		return nil
	}

	// check the value of the arguments struct
	aVal := reflect.ValueOf(args)

	for _, fIndex := range argumentsFilterIndexes {
		v := aVal.FieldByIndex(fIndex)

		if !v.IsZero() { // flag was set iff it is non-zero
			tp := argumentsType.FieldByIndex(fIndex) // needed for the error message only!
			return errTakesNoArgument.WithMessageF(args.Command, "--"+tp.Tag.Get("long"))
		}
	}

	return nil
}

// AllowsOption checks if the provided option is allowed by this option
func (req Requirement) AllowsOption(opt usagefmt.Opt) bool {
	return req.AllowsFilter || text.SliceContainsAny(opt.Long(), argumentsGeneralOptions...)
}

// reflect access to the arguments type
// TODO: Copied over from program
var argumentsType reflect.Type = reflect.TypeOf((*program.Arguments)(nil)).Elem() // TypeOf[Arguments]

var argumentsGeneralOptions []string // names of options that are considered non-filter
var argumentsFilterIndexes [][]int   // indexes of filter options

func init() {
	// iterate over the fields of the type
	fieldCount := argumentsType.NumField()
	for i := 0; i < fieldCount; i++ {
		field := argumentsType.Field(i)

		// skip over options that do not have a 'long' name
		longName, hasLongName := field.Tag.Lookup("long")
		if !hasLongName {
			continue
		}

		// argument is a nonfilter argument!
		if field.Tag.Get("nofilter") == "true" {
			argumentsGeneralOptions = append(argumentsGeneralOptions, longName)
			continue
		}

		// it's a long filter name
		argumentsFilterIndexes = append(argumentsFilterIndexes, field.Index)
	}
}
