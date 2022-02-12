package usagefmt

import (
	"strings"

	"github.com/jessevdk/go-flags"
)

// Opt represents an option
type Opt struct {

	// For the purposes of documentation we use the following argument as an example.
	//   -n, --number digit  A digit used within something (default: 42)

	Required bool // false

	Short     []string // ["n"]
	Long      []string // ["number"]
	FieldName string   // name of the underlying field name

	Value string // "digit"

	Usage   string // "A digit used within something"
	Default string // "42"
}

// NewOpt returns a new Arg based on an Option
func NewOpt(option *flags.Option) (opt Opt) {
	opt.Required = option.Required

	short := option.ShortName
	if short != rune(0) {
		opt.Short = []string{string(short)}
	}

	long := option.LongName
	if long != "" {
		opt.Long = []string{long}
	}

	opt.FieldName = option.Field().Name

	opt.Value = option.ValueName

	opt.Usage = option.Description

	dflt := option.Default
	if len(dflt) != 0 {
		opt.Default = strings.Join(dflt, ", ")
	}

	return
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
