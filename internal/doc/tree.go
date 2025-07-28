package doc

import (
	"iter"

	"github.com/spf13/cobra"
)

// AllCommands returns a sequence of all commands in the command tree.
func AllCommands(cmd *cobra.Command) iter.Seq[*cobra.Command] {
	return func(yield func(*cobra.Command) bool) {
		allCommands(cmd, yield)
	}
}

func allCommands(cmd *cobra.Command, yield func(*cobra.Command) bool) bool {
	if !yield(cmd) {
		return false
	}
	for _, c := range cmd.Commands() {
		if !allCommands(c, yield) {
			return false
		}
	}
	return true
}

// CountCommands returns the number of commands in the command tree.
func CountCommands(cmd *cobra.Command) (count int) {
	for range AllCommands(cmd) {
		count++
	}
	return count
}
