package cmd

//spellchecker:words github cobra ggman
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
)

// TODO: move any error declarations here

func NewTODOCommand() *cobra.Command {
	// TODO: introduce new var s for flags here

	todo := &cobra.Command{
		Use:   "todo",       // todo: put in old command name
		Short: "",           // todo: copy old description
		Args:  cobra.NoArgs, // todo: pick the appropriate thing from cobra

		PreRunE: func(cmd *cobra.Command, args []string) error {
			ggman.SetRequirements(cmd, &env.Requirement{
				// TODO: set correctly from the old Description method
				// set the actual flags, don't make a reference to the old struct
			})

			// TODO: do anything the old AfterParse() method does
			// if not, just leave as is.
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := ggman.GetEnv(cmd)
			if err != nil {
				return fmt.Errorf("failed to get environment: %w", err)
			}

			// TODO: implement command here - what the old Run method did.
			// any methods on the old context.Environment should use the env variable above.
			// if you need to pass an stream.IOStream somewhere, use the [streamFromCommand] method.
			_ = env
			return nil
		},
	}

	flags := todo.Flags()
	// TODO: put parse flag stuff here
	// extract from the contents of the old struct and instead create appropriate vars
	_ = flags

	return todo
}

// TODO: this struct shouldn't be referenced anymore once you're done with it.
