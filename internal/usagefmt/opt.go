package usagefmt

import (
	"strings"

	"github.com/jessevdk/go-flags"
)

// Opt is anything that can be treated as an option for a usage message.
type Opt interface {

	// For the purposes of documentation we use the following argument as an example.
	//   -n, --number digit  A digit used within something (default: 42)

	Required() bool // false

	Short() []string // ["n"]
	Long() []string  // ["number"]

	Value() string // "digit"

	Usage() string   // "A digit used within something"
	Default() string // "42"
}

// NewOpt returns a new Arg based on an Option
func NewOpt(option *flags.Option) Opt {
	return opt{Option: option}
}

// opt implements Option
type opt struct {
	*flags.Option
}

// Required checks if this argument is required
func (o opt) Required() bool {
	return o.Option.Required
}

// Short returns a list of short names
func (o opt) Short() []string {
	short := o.Option.ShortName
	if short == rune(0) {
		return nil
	}
	return []string{string(short)}
}

// Long returns a list of long names
func (o opt) Long() []string {
	long := o.Option.LongName
	if long == "" {
		return nil
	}
	return []string{long}
}

// Value returns the name of the argument
func (o opt) Value() string {
	return o.Option.ValueName
}

// Usage returns a human-readable usage text
func (o opt) Usage() string {
	return o.Option.Description
}

// Default returns the default value of this argument as a string
// or the empty string if there is no default.
func (o opt) Default() string {
	dflt := o.Option.Default
	if len(dflt) == 0 {
		return ""
	}
	return strings.Join(dflt, ", ")
}

// MakeOpts returns all options available within parser.
// See also NewOpt.
//
// This function is untested.
func MakeOpts(parser *flags.Parser) (opts []Opt) {

	// collect all the options
	var options []*flags.Option
	groups := parser.Groups()
	for _, g := range groups {
		options = append(options, g.Options()...)
	}

	// make them into proper opts
	opts = make([]Opt, len(options))
	for i, opt := range options {
		opts[i] = NewOpt(opt)
	}

	return
}
