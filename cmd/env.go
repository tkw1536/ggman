package cmd

//spellchecker:words essio shellescape github cobra ggman pkglib collection exit
import (
	"fmt"

	"al.essio.dev/pkg/shellescape"
	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/collection"
	"go.tkw01536.de/pkglib/exit"
)

func NewEnvCommand() *cobra.Command {
	impl := new(_env)

	cmd := &cobra.Command{
		Use:   "env [VAR...]",
		Short: "print information about the ggman environment",
		Long: `Env prints "name=value" pairs about the environment the ggman command is running in to standard output.
value is escaped for use in a shell.

By default, env prints information about all known variables.
To print information about a subset of variables, they can be provided as positional arguments.
Variables names are matched case-insensitively.`,
		Args: cobra.ArbitraryArgs,

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.List, "list", "l", false, "instead of \"name=value\" pairs print only the variable")
	flags.BoolVarP(&impl.Describe, "describe", "d", false, "instead of \"name=value\" pairs print \"name: description\" pairs describing the use of variables")
	flags.BoolVarP(&impl.Raw, "raw", "r", false, "instead of \"name=value\" pairs print only the unescaped value")

	return cmd
}

type _env struct {
	Positionals struct {
		Vars []string
	}

	List     bool
	Describe bool
	Raw      bool
}

var (
	errEnvInvalidVar        = exit.NewErrorWithCode("unknown environment variable", exit.ExitCommandArguments)
	errEnvModesIncompatible = exit.NewErrorWithCode("at most one of `--raw`, `--list` and `--describe` may be given", exit.ExitCommandArguments)
)

func (e *_env) ParseArgs(cmd *cobra.Command, args []string) error {
	// check that at most one mode was given
	count := 0
	if e.Describe {
		count++
	}
	if e.Raw {
		count++
	}
	if e.List {
		count++
	}
	if count > 1 {
		return errEnvModesIncompatible
	}

	e.Positionals.Vars = args
	return nil
}

func (e *_env) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{NeedsRoot: true})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	variables, err := e.variables()
	if err != nil {
		return err
	}

	for _, v := range variables {
		switch {
		case e.List:
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), v.Key); err != nil {
				return fmt.Errorf("%w: %w", errGenericOutput, err)
			}
		case e.Raw:
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), v.Get(environment)); err != nil {
				return fmt.Errorf("%w: %w", errGenericOutput, err)
			}
		case e.Describe:
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", v.Key, v.Description); err != nil {
				return fmt.Errorf("%w: %w", errGenericOutput, err)
			}
		default:
			value := shellescape.Quote(v.Get(environment))
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", v.Key, value); err != nil {
				return fmt.Errorf("%w: %w", errGenericOutput, err)
			}
		}
	}

	return nil
}

func (e *_env) variables() ([]env.UserVariable, error) {
	// no variables provided => use all of them
	if len(e.Positionals.Vars) == 0 {
		return env.GetUserVariables(), nil
	}

	var invalid string
	variables := collection.MapSlice(e.Positionals.Vars, func(name string) env.UserVariable {
		value, ok := env.GetUserVariable(name)
		if !ok && invalid == "" { // store an invalid name!
			invalid = name
		}
		return value
	})

	if invalid != "" {
		return nil, fmt.Errorf("%q: %w", invalid, errEnvInvalidVar)
	}

	return variables, nil
}
