package usagefmt

import (
	"strings"
)

// Page represents data required about a command to generate a help page
type Page struct {
	// Name of the main executable (e.g. "ggman") and global options
	MainName string
	MainOpts []Opt

	// a general description of the command or subcommand
	Description string

	// Name of available subcommands, only used for general help page
	SubCommands []string

	// Name of the subcommand (e.g. "ls" for "ggman ls") and specific options
	// omit SubName (or leave it empty) to generate a help page without a command
	SubName string
	SubOpts []Opt

	// named argument specification, see FmtSpecName.
	MetaName         string
	MetaMin, MetaMax int

	// Usage description
	Usage string
}

// String generates a help page for this command
func (page Page) String() string {
	// grab a builder from the pool
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	page.Build(builder)
	return builder.String()
}

// Build writes this page into builder
func (page Page) Build(builder *strings.Builder) {
	if page.SubName != "" {
		page.sub(builder)
		return
	}
	page.main(builder)
}

// subSpec is spec for a subcommand
const subSpec = "COMMAND [ARGS...]"

// subMsgTpl is the usage message of a subcommand.
// It consists of two parts.
const (
	// subMsgTpl = subMsg1 + "%s" + subMsg2
	subMsg1 = "Command to call. One of "
	subMsg2 = ". See individual commands for more help."
)

// main implements generating a help page for a main command (without a subcommand)
func (page Page) main(builder *strings.Builder) {

	//
	// Command specification
	//

	// main command
	builder.WriteString("Usage: ")
	builder.WriteString(page.MainName)

	for _, arg := range page.MainOpts {
		builder.WriteRune(' ')
		SpecShort(builder, arg)
	}

	builder.WriteString(" [--] ")
	builder.WriteString(subSpec)

	// description (if any)
	if page.Description != "" {
		builder.WriteString("\n\n")
		builder.WriteString(page.Description)
	}

	//
	// Argument description
	//

	for _, arg := range page.MainOpts {
		Message(builder, arg)
	}

	// write a usage message for the commands

	builder.WriteString(usageMsg1)
	builder.WriteString(subSpec)
	builder.WriteString(usageMsg2)

	// replace the list of commands in subMsgTpl
	builder.WriteString(subMsg1)
	Commands(builder, page.SubCommands)
	builder.WriteString(subMsg2)

	builder.WriteString(usageMsg3)
}

// sub implements generating a help page for a non-empty subcommand.
func (page Page) sub(builder *strings.Builder) {

	//
	// Command specification
	//

	// main command
	builder.WriteString("Usage: ")
	builder.WriteString(page.MainName)

	for _, arg := range page.MainOpts {
		builder.WriteRune(' ')
		SpecShort(builder, arg)
	}

	if len(page.MainOpts) >= 0 {
		builder.WriteString(" [--]")
	}

	// subcommand
	builder.WriteRune(' ')
	builder.WriteString(page.SubName)

	for _, arg := range page.SubOpts {
		builder.WriteRune(' ')
		SpecShort(builder, arg)
	}

	if page.MetaMax != 0 {
		builder.WriteString(" [--] ")
		SpecPositional(builder, page.MetaName, page.MetaMin, page.MetaMax)
	}

	// description (if any)
	if page.Description != "" {
		builder.WriteString("\n\n")
		builder.WriteString(page.Description)
	}

	//
	// Argument description
	//

	builder.WriteString("\n\nGlobal Arguments:")
	for _, opt := range page.MainOpts {
		Message(builder, opt)
	}

	// no command arguments provided!
	if len(page.SubOpts) == 0 && page.MetaMax == 0 {
		return
	}

	builder.WriteString("\n\nCommand Arguments:")

	for _, opt := range page.SubOpts {
		Message(builder, opt)
	}

	// write the usage message of the argument (if any)
	if page.Usage != "" {
		builder.WriteString(usageMsg1)
		SpecPositional(builder, page.MetaName, page.MetaMin, page.MetaMax)
		builder.WriteString(usageMsg2)
		builder.WriteString(page.Usage)
		builder.WriteString(usageMsg3)
	}
}
