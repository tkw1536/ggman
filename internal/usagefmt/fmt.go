// Package usagefmt provides facilities for formatting usage messages.
package usagefmt

import (
	"strconv"
	"strings"
	"sync"

	"github.com/tkw1536/ggman/internal/text"
)

// builderPool used by various formatters in this package
var builderPool = &sync.Pool{
	New: func() interface{} { return new(strings.Builder) },
}

// SpecShort writes a short specification of the option into builder.
// It is of the form '--flag|-f value'.
// SpecShort adds braces around the argument if it is optional.
func SpecShort(builder *strings.Builder, opt Opt) {
	spec(builder, opt, "|", true, true)
}

// SpecLong writes a long specification of the option into builder.
// It is of the form '-f, --flag value'.
func SpecLong(builder *strings.Builder, opt Opt) {
	spec(builder, opt, ", ", false, false)
}

// FmtSpecShort is like SpecShort, but returns a string
func FmtSpecShort(opt Opt) string {
	// grab a builder from the pool
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	SpecShort(builder, opt)
	return builder.String()
}

// FmtSpecLong is like SpecLong, but returns a string
func FmtSpecLong(opt Opt) string {
	// grab a builder from the pool
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	SpecLong(builder, opt)
	return builder.String()
}

// spec implements SpecShort and SpecLong.
//
// sep indicates how to seperate arguments.
// longFirst indicates that long argument names should be listed before short arguments.
// optionalBraces indicates if braces should be placed around the argument if it is optional.
func spec(builder *strings.Builder, opt Opt, sep string, longFirst bool, optionalBraces bool) {
	// if the argument is optional put braces around it!
	if optionalBraces && !opt.Required() {
		builder.WriteString("[")
		defer builder.WriteString("]")
	}

	// collect long and short arguments and combine them
	la := text.SliceCopy(opt.Long())
	for k, v := range la {
		la[k] = "--" + v
	}

	sa := text.SliceCopy(opt.Short())
	for k, v := range sa {
		sa[k] = "-" + v
	}

	// write the joined versions of the arguments into the specification
	var args []string
	if longFirst {
		args = append(la, sa...)
	} else {
		args = append(sa, la...)
	}
	text.Join(builder, args, sep)

	// write the value (if any)
	if value := opt.Value(); value != "" {
		builder.WriteRune(' ')
		builder.WriteString(value)
	}
}

// usageMsgTpl is the template for long usage messages
// it is split into three parts, that are joined by the arguments.
//
//  const usageMsgTpl = usageMsg1 + "%s" + usageMsg2 + "%s" + usageMsg3
const (
	usageMsg1 = "\n\n   "
	usageMsg2 = "\n      "
	usageMsg3 = ""
)

// Message writes a long message describing the argument to builder
// It is of the form '-f, --flag ARG' and 'DESCRIPTION (default DEFAULT)'
//
// This function is implicity tested via other tests.
func Message(builder *strings.Builder, opt Opt) {

	builder.WriteString(usageMsg1)
	SpecLong(builder, opt)
	builder.WriteString(usageMsg2)

	builder.WriteString(opt.Usage())
	if dflt := opt.Default(); dflt != "" {
		builder.WriteString(" (default ")
		builder.WriteString(dflt)
		builder.WriteString(")")
	}

	builder.WriteString(usageMsg3)
}

// FmtMessage is like Message, but returns a string
func FmtMessage(opt Opt) string {
	// grab a builder from the pool
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	Message(builder, opt)
	return builder.String()
}

// DefaultPositionalName is the default name used for a positional argument.
// See SpecPositional.
const DefaultPositionalName = "ARGUMENT"

// SpecPositional creates a spec for a positional argument e.g. "arg [arg...]" and writes it to builder.
//
// name is the name of the named argument, min and max are the minimum and maximum respectively.
// when name is the empty string, uses DefaultPositional.
//
// min must be non-negative. max must be bigger than min or less than 0.
// when max is 0, assumes that the argument can be repeated indefinitly.
func SpecPositional(builder *strings.Builder, name string, min, max int) {
	extra := max - min // extra is the number of optional argument

	if min < 0 || (max > 0 && extra < 0) { // invalid arguments
		panic("NameSpec: negative min or out of range max")
	}

	if name == "" {
		name = DefaultPositionalName
	}

	// nothing to generate!
	if max == 0 && extra == 0 {
		return
	}

	// arg arg arg
	text.RepeatJoin(builder, name, " ", min)
	if min > 0 && extra != 0 {
		builder.WriteString(" ")
	}

	if max < 0 {
		// [arg ...]
		builder.WriteString("[")
		builder.WriteString(name)
		builder.WriteString(" ...]")
		return
	}

	// [arg [arg]]
	text.RepeatJoin(builder, "["+name, " ", extra)
	text.Repeat(builder, "]", extra)
}

// FmtSpecPositional is like SpecPositional except that it returns a string.
func FmtSpecPositional(name string, min, max int) string {
	// grab a builder from the pool
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	SpecPositional(builder, name, min, max)
	return builder.String()
}

// Commands writes a human readable representation of commands into builder.
func Commands(builder *strings.Builder, commands []string) {
	if len(commands) == 0 {
		return
	}
	builder.WriteString(strconv.Quote(commands[0]))
	for _, cmd := range commands[1:] {
		builder.WriteString(", ")
		builder.WriteString(strconv.Quote(cmd))
	}
}

// FmtCommands is like Commands, but returns a string.
func FmtCommands(commands []string) string {
	// grab a builder from the pool
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	Commands(builder, commands)
	return builder.String()
}
