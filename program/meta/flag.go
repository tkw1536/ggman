package meta

import (
	"io"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/internal/slice"
	"github.com/tkw1536/ggman/internal/text"
)

// Flag holds meta-information about a single flag of a command.
// It is similar to "github.com/jessevdk/go-flags".Flag.
//
// Unlike the actual flag, it holds no reflect references, and does not contain an actual value.
//
// See also NewFlag.
type Flag struct {

	// For the purposes of documentation we use the following argument as an example.
	//   -n, --number digit  A digit used within something (default: 42)

	// The name of the underlying struct field this flag comes from.
	FieldName string // "Number"

	// Short and Long Names of the flag
	// each potentially more than one
	Short []string // ["n"]
	Long  []string // ["number"]

	// Indicates if the flag is required
	Required bool // false

	// Name and Description of the flag in help texts
	Value string // "digit"
	Usage string // "A digit used within something"

	// Default value of the flag (as a string)
	Default string // "42"
}

// NewFlag creates a new flag based on an option from the flags package.
func NewFlag(option *flags.Option) (flag Flag) {
	flag.Required = option.Required

	short := option.ShortName
	if short != rune(0) {
		flag.Short = []string{string(short)}
	}

	long := option.LongName
	if long != "" {
		flag.Long = []string{long}
	}

	flag.FieldName = option.Field().Name

	flag.Value = option.ValueName

	flag.Usage = option.Description

	dflt := option.Default
	if len(dflt) != 0 {
		flag.Default = strings.Join(dflt, ", ")
	}

	return
}

// AllFlags returns all flags available to parser.
//
// This function is untested.
func AllFlags(parser *flags.Parser) []Flag {
	// collect all the options
	var options []*flags.Option
	groups := parser.Groups()
	for _, g := range groups {
		options = append(options, g.Options()...)
	}

	// turn them into proper flags
	flags := make([]Flag, len(options))
	for i, opt := range options {
		flags[i] = NewFlag(opt)
	}
	return flags
}

// WriteSpecTo writes a short specification of f into w.
// It is of the form
//    --flag|-f value
// WriteSpecTo adds braces around the argument if it is optional.
func (f Flag) WriteSpecTo(w io.Writer) {
	f.spec(w, "|", true, true)
}

// WriteLongSpecTo writes a long specification of f into w.
// It is of the form
//  -f, --flag value
// WriteLongSpecTo does not add any brackets around the argument.
func (opt Flag) WriteLongSpecTo(w io.Writer) {
	opt.spec(w, ", ", false, false)
}

// spec implements SpecShort and SpecLong.
//
// sep indicates how to seperate arguments.
// longFirst indicates that long argument names should be listed before short arguments.
// optionalBraces indicates if braces should be placed around the argument if it is optional.
func (opt Flag) spec(w io.Writer, sep string, longFirst bool, optionalBraces bool) {
	// if the argument is optional put braces around it!
	if optionalBraces && !opt.Required {
		io.WriteString(w, "[")
		defer io.WriteString(w, "]")
	}

	// collect long and short arguments and combine them
	la := slice.Copy(opt.Long)
	for k, v := range la {
		la[k] = "--" + v
	}

	sa := slice.Copy(opt.Short)
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
	text.Join(w, args, sep)

	// write the value (if any)
	if value := opt.Value; value != "" {
		io.WriteString(w, " ")
		io.WriteString(w, value)
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

// WriteMessageTo writes a long message of f to w.
// It is of the form
//    -f, --flag ARG
// and
//    DESCRIPTION (default DEFAULT)
// .
//
// This function is implicity tested via other tests.
func (opt Flag) WriteMessageTo(w io.Writer) {

	io.WriteString(w, usageMsg1)
	opt.WriteLongSpecTo(w)
	io.WriteString(w, usageMsg2)

	io.WriteString(w, opt.Usage)
	if dflt := opt.Default; dflt != "" {
		io.WriteString(w, " (default ")
		io.WriteString(w, dflt)
		io.WriteString(w, ")")
	}

	io.WriteString(w, usageMsg3)
}
