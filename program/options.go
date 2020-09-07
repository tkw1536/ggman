package program

import (
	"fmt"
	"reflect"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/tkw1536/ggman/env"
)

// Options represent the options for a specific command
type Options struct {
	Environment env.Requirement

	// minimum and maximum number of arguments
	MinArgs int
	MaxArgs int

	// the name of the metavar to use for the usage string
	Metavar string

	// Description of the argument
	UsageDescription string

	// Description of the flag
	FlagDescription string
}

// usageTemplate is a template string used to display usage
const usageTemplate = "\n\n   %s\n       %s"

// Usage returns a string representing the usage of the options induced by this command
func (opt Options) Usage(name string, flagset *flag.FlagSet) (usage string) {

	// gather all the usage information of all the flags
	flags := make([]flagUsage, 0)
	flagset.VisitAll(func(f *flag.Flag) {
		flags = append(flags, flagUsage{flag: f})
	})

	usage = "Usage: ggman"

	// the for argument
	if opt.Environment.AllowsFilter {
		usage += " [for|--for|-f FILTER]"
	}

	// the name and help
	usage += " " + name + " [help|--help|-h]"

	// read the metavar
	mv := opt.Metavar
	if mv == "" {
		mv = "ARGUMENT"
	}

	for _, f := range flags {
		usage += fmt.Sprintf(" [%s]", f.Flag())
	}

	usage += " [--]"

	argSyntax := ""
	if opt.MaxArgs == -1 {
		// write out the argument an appropriate number of times
		argSyntax += strings.Repeat(" "+mv, opt.MinArgs)
		argSyntax += " [" + mv + " ... ]"
	} else {
		// write out the argument an appropriate number of times
		argSyntax += strings.Repeat(" "+mv, opt.MinArgs)
		argSyntax += strings.Repeat(" ["+mv, opt.MaxArgs-opt.MinArgs)
		argSyntax += strings.Repeat("]", opt.MaxArgs-opt.MinArgs)
	}

	usage += argSyntax

	// start with the help argument
	usage += fmt.Sprintf(usageTemplate, "help|--help|-h", "Print this usage message and exit.")

	// contineu with the 'for' argument
	if opt.Environment.AllowsFilter {
		usage += fmt.Sprintf(usageTemplate, "for|--for|-f FILTER", "Filter the list of repositories to apply command to by FILTER.")
	}

	// add description for the flagset
	for _, u := range flags {
		usage += u.Description()
	}

	// add the description of the arguments
	if opt.UsageDescription != "" {
		usage += fmt.Sprintf(usageTemplate, strings.TrimLeft(argSyntax, " "), opt.UsageDescription)
	}

	return
}

// flagUsage is a utility struct for descriptions about flags
// It has been adapted from the defaultUsage funtions in the flag package.
type flagUsage struct {
	flag *flag.Flag
}

// Flag returns the flag in the form '--flag|-f value'
func (f flagUsage) Flag() string {
	s := fmt.Sprintf("--%s", f.flag.Name)
	if shorthand := f.flag.Shorthand; shorthand != "" {
		s += fmt.Sprintf("|-%s", shorthand)
	}

	if name, _ := flag.UnquoteUsage(f.flag); name != "" {
		s += " " + name
	}
	return s
}

// Description returns a long form description of the flag.
// It is of the form '-f, --flag ARG' and 'DESCRIPTION (default DEFAULT)'
//
// This function has been adapated from flag.PrintDefaults
func (f flagUsage) Description() string {
	// extract the short flag de
	description := fmt.Sprintf("--%s", f.flag.Name)
	if shorthand := f.flag.Shorthand; shorthand != "" {
		description = fmt.Sprintf("-%s, %s", shorthand, description)
	}
	name, usage := flag.UnquoteUsage(f.flag)
	if len(name) > 0 {
		description += " " + name
	}

	// build the usage string, and add the default
	usage = strings.ReplaceAll(usage, "\n", "")
	if !f.isZeroValue() {
		usage += fmt.Sprintf(" (default %v)", f.flag.DefValue)
	}

	return fmt.Sprintf(usageTemplate, description, usage)
}

// isZeroValue determines whether the default string represents the zero
// value for a flag. Adapated from flag "github.com/spf13/pflag".isZeroValue.
func (f flagUsage) isZeroValue() bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(f.flag.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	return f.flag.DefValue == z.Interface().(flag.Value).String()
}
