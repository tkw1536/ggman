// Package cmd implements 'ggman' subcommands.
//
// Each subcommand is represented as a single variable of type program.Command.
package cmd

import "github.com/tkw1536/ggman/program"

// All returns a channel which will be sent each command.
// Once all commands have been sent, the channel will be closed.
// The order of commands is not undefined.
//
// This function is untested.
func All() <-chan program.Command {
	cmds := make(chan program.Command)

	go func() {
		defer func() { close(cmds) }()

		cmds <- Root

		cmds <- Ls
		cmds <- Lsr

		cmds <- Where

		cmds <- Canon

		cmds <- Comps

		cmds <- Fetch
		cmds <- Pull

		cmds <- Fix

		cmds <- Clone

		cmds <- Link

		cmds <- License

		cmds <- Here

		cmds <- Web
		cmds <- URL

		cmds <- FindBranch

		cmds <- Relocate
	}()

	return cmds
}
