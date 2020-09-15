package program

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/pflag"
	"github.com/tkw1536/ggman/constants"
)

// flagUsage is a utility struct for descriptions about flags
// It has been adapted from the defaultUsage functions in the flag package.
type flagUsage struct {
	flag *pflag.Flag
}

// Flag returns the flag in the form '--flag|-f value'
func (f flagUsage) Flag() string {
	s := fmt.Sprintf("--%s", f.flag.Name)
	if shorthand := f.flag.Shorthand; shorthand != "" {
		s += fmt.Sprintf("|-%s", shorthand)
	}

	if name, _ := pflag.UnquoteUsage(f.flag); name != "" {
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
	name, usage := pflag.UnquoteUsage(f.flag)
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
// value for a flag. Adapated from flag.isZeroValue.
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
	return f.flag.DefValue == z.Interface().(pflag.Value).String()
}

// usageTemplate is a template string used to display usage
const usageTemplate = "\n\n   %s\n       %s"

//
// PROGRAM
//

// Usage returns a string describing the usage of this flagset
func (p Program) Usage(flagset *pflag.FlagSet) (usage string) {

	usage = fmt.Sprintf("ggman version %s\n\n", constants.BuildVersion)

	// gather all the usage information of all the flags
	flags := make([]flagUsage, 0)
	flagset.VisitAll(func(f *pflag.Flag) {
		flags = append(flags, flagUsage{flag: f})
	})

	usage += "Usage: ggman"

	for _, f := range flags {
		usage += fmt.Sprintf(" [%s]", f.Flag())
	}

	usage += " [--] COMMAND [ARGS...]"

	// add description for the flagset
	for _, u := range flags {
		usage += u.Description()
	}

	usage += fmt.Sprintf(usageTemplate, "COMMAND [ARGS...]", fmt.Sprintf("Command to call. One of %s. See individual commands for more help.", p.knownCommandsString()))

	usage += "\n\nggman is licensed under the terms of the MIT License. Use 'ggman license' to view licensing information."

	return usage
}

// knownCommandsString returns a string containing all the known commands
func (p Program) knownCommandsString() string {
	keys := make([]string, 0, len(p.commands))
	for k := range p.commands {
		keys = append(keys, "'"+k+"'")
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

//
// COMMAND
//

// Usage returns a string representing the usage of the options induced by this command
func (opt Options) Usage(name string, flagset *pflag.FlagSet) (usage string) {

	// gather all the usage information of all the flags
	flags := make([]flagUsage, 0)
	flagset.VisitAll(func(f *pflag.Flag) {
		flags = append(flags, flagUsage{flag: f})
	})

	usage = "Usage: ggman"

	// the for argument
	if opt.Environment.AllowsFilter {
		usage += " [--for|-f FILTER]"
	}

	usage += " [global arguments] [--]"

	// the name and help
	usage += " " + name + " [--help|-h]"

	// read the metavar
	mv := opt.Metavar
	if mv == "" {
		mv = "ARGUMENT"
	}

	for _, f := range flags {
		usage += fmt.Sprintf(" [%s]", f.Flag())
	}

	if opt.MinArgs > 0 {
		usage += " [--]"
	}

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
	usage += fmt.Sprintf(usageTemplate, "-h, --help", "Print this usage message and exit.")

	// contineu with the 'for' argument
	if opt.Environment.AllowsFilter {
		usage += fmt.Sprintf(usageTemplate, "-f, --for filter", "Filter the list of repositories to apply command to by FILTER.")
	}

	usage += fmt.Sprintf(usageTemplate, "global arguments", "Global arguments for ggman. See ggman --help for more information.")

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
