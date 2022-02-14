package meta

import (
	"io"
	"strconv"
)

// Meta holds meta-information about an entire program or a subcommand.
// It is used to generate a usage page.
type Meta struct {
	// Name of the Executable and Current command.
	// When Command is empty, the entire struct describes the program as a whole.
	Executable string
	Command    string

	// Description holds a human-readable description of the object being described.
	Description string

	// Applicable Global, Command and Positional Flags.
	GlobalFlags  []Flag
	CommandFlags []Flag
	Positional   Positional

	// List of available sub-commands, only set when Command == "".
	Commands []string
}

// WriteMessageTo writes the human-readable message of this meta into w
func (meta Meta) WriteMessageTo(w io.Writer) {
	if meta.Command != "" {
		meta.writeCommandMessageTo(w)
		return
	}
	meta.writeProgramMessageTo(w)
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

func (meta Meta) writeProgramMessageTo(w io.Writer) {
	//
	// Command specification
	//

	// main command
	io.WriteString(w, "Usage: ")
	io.WriteString(w, meta.Executable)

	for _, arg := range meta.GlobalFlags {
		io.WriteString(w, " ")
		arg.WriteSpecTo(w)
	}

	io.WriteString(w, " [--] ")
	io.WriteString(w, subSpec)

	// description (if any)
	if meta.Description != "" {
		io.WriteString(w, "\n\n")
		io.WriteString(w, meta.Description)
	}

	//
	// Argument description
	//

	for _, arg := range meta.GlobalFlags {
		arg.WriteMessageTo(w)
	}

	// write a usage message for the commands

	io.WriteString(w, usageMsg1)
	io.WriteString(w, subSpec)
	io.WriteString(w, usageMsg2)

	// replace the list of commands in subMsgTpl
	io.WriteString(w, subMsg1)
	meta.writeCommandsTo(w)
	io.WriteString(w, subMsg2)

	io.WriteString(w, usageMsg3)
}

// WriteCommandsTo writes the list of commands to w.
func (meta Meta) writeCommandsTo(w io.Writer) {
	if len(meta.Commands) == 0 {
		return
	}
	io.WriteString(w, strconv.Quote(meta.Commands[0]))
	for _, cmd := range meta.Commands[1:] {
		io.WriteString(w, ", ")
		io.WriteString(w, strconv.Quote(cmd))
	}
}

func (page Meta) writeCommandMessageTo(w io.Writer) {

	//
	// Command specification
	//

	// main command
	io.WriteString(w, "Usage: ")
	io.WriteString(w, page.Executable)

	for _, arg := range page.GlobalFlags {
		io.WriteString(w, " ")
		arg.WriteSpecTo(w)
	}

	if len(page.GlobalFlags) >= 0 {
		io.WriteString(w, " [--]")
	}

	// subcommand
	io.WriteString(w, " ")
	io.WriteString(w, page.Command)

	for _, arg := range page.CommandFlags {
		io.WriteString(w, " ")
		arg.WriteSpecTo(w)
	}

	if page.Positional.Max != 0 {
		io.WriteString(w, " [--] ")
		page.Positional.WriteSpecTo(w)
	}

	// description (if any)
	if page.Description != "" {
		io.WriteString(w, "\n\n")
		io.WriteString(w, page.Description)
	}

	//
	// Argument description
	//

	io.WriteString(w, "\n\nGlobal Arguments:")
	for _, opt := range page.GlobalFlags {
		opt.WriteMessageTo(w)
	}

	// no command arguments provided!
	if len(page.CommandFlags) == 0 && page.Positional.Max == 0 {
		return
	}

	io.WriteString(w, "\n\nCommand Arguments:")

	for _, opt := range page.CommandFlags {
		opt.WriteMessageTo(w)
	}

	// write the usage message of the argument (if any)
	if page.Positional.Description != "" {
		io.WriteString(w, usageMsg1)
		page.Positional.WriteSpecTo(w)
		io.WriteString(w, usageMsg2)
		io.WriteString(w, page.Positional.Description)
		io.WriteString(w, usageMsg3)
	}
}
