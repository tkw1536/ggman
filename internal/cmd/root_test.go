package cmd_test

//spellchecker:words slices strings testing github cobra ggman internal mockenv
import (
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/ggman/internal/mockenv"
)

// This test invokes all commands with the --help flag.
func Test_NewCommand(t *testing.T) {
	t.Parallel()

	// find all the commands that have flag parsing enabled
	root := cmd.NewCommand(t.Context(), env.Parameters{})
	commands := findCommands(root, func(cmd *cobra.Command) bool {
		return !cmd.DisableFlagParsing
	})

	// run each of the commands with the --help flag
	for _, argv := range commands {
		t.Run(strings.Join(argv, " "), func(t *testing.T) {
			t.Parallel()

			mock := mockenv.NewMockEnv(t)

			code, _, _ := mock.Run(t, nil, cmd.NewCommand, "", "", append(argv[1:], "--help")...)
			if code != 0 {
				t.Errorf("command %q failed with code %d", strings.Join(argv, " "), code)
			}
		})
	}
}

func findCommands(cmd *cobra.Command, include func(cmd *cobra.Command) bool) (commands [][]string) {
	var walkTree func(cmd *cobra.Command, path []string)
	walkTree = func(cmd *cobra.Command, path []string) {
		us := slices.Clone(path)
		us = append(us, cmd.Name())

		if include(cmd) {
			commands = append(commands, us)
		}

		for _, sub := range cmd.Commands() {
			walkTree(sub, us)
		}
	}
	walkTree(cmd, nil)
	return
}
