package cmd

//spellchecker:words encoding json jsontext github cobra ggman internal pkglib exit
import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words wrapcheck

func NewParseCommand() *cobra.Command {
	impl := new(parse)

	cmd := &cobra.Command{
		Use:   "parse URL",
		Short: "Parse a URL into components and forge tree reference",
		Long: `Parse prints components and forge tree reference of a URL.

By default output is JSON with fields 'comps', 'branch' and 'path'.

    ggman parse https://github.com/hello/world/tree/main/src

The '--comps', '--branch' and '--path' flags print only that field as plain text.
'--comps' prints one component per line.

    ggman parse --comps https://github.com/hello/world.git

The '--no-forge-split' flag skips forge tree splitting.

'ggman comps' is an alias for 'ggman parse --comps'.`,
		Args: cobra.ExactArgs(1),

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.Comps, "comps", false, "print only URL components, one per line")
	flags.BoolVar(&impl.Branch, "branch", false, "print only the forge tree branch/ref")
	flags.BoolVar(&impl.Path, "path", false, "print only the relative path within the forge tree")
	flags.BoolVar(&impl.NoForgeSplit, "no-forge-split", false, "do not split forge tree references from the URL")

	return cmd
}

type parse struct {
	Positional struct {
		URL    env.URL
		Branch string
		Path   string
	}
	Comps        bool
	Branch       bool
	Path         bool
	NoForgeSplit bool
}

type parseResult struct {
	Comps  []string `json:"comps"`
	Branch string   `json:"branch"`
	Path   string   `json:"path"`
}

var (
	errParseExclusiveFlags = exit.NewErrorWithCode(`only one of "--comps", "--branch" and "--path" may be provided`, env.ExitCommandArguments)
)

func (p *parse) ParseArgs(cmd *cobra.Command, args []string) error {
	count := 0
	if p.Comps {
		count++
	}
	if p.Branch {
		count++
	}
	if p.Path {
		count++
	}
	if count > 1 {
		return errParseExclusiveFlags
	}

	var url env.URL
	var ref, relative string
	if p.NoForgeSplit {
		url = env.ParseURL(args[0])
	} else {
		url, _, ref, relative = env.ParseURLAndForgeReference(args[0])
	}

	p.Positional.URL = url
	p.Positional.Branch = ref
	p.Positional.Path = relative
	return nil
}

func (p *parse) Exec(cmd *cobra.Command, args []string) error {
	result := parseResult{
		Comps:  p.Positional.URL.Components(),
		Branch: p.Positional.Branch,
		Path:   p.Positional.Path,
	}
	if result.Comps == nil {
		result.Comps = []string{}
	}

	switch {
	case p.Comps:
		for _, comp := range result.Comps {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), comp); err != nil {
				return fmt.Errorf("%w: %w", errGenericOutput, err)
			}
		}
	case p.Branch:
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), result.Branch); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
	case p.Path:
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), result.Path); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
	default:
		if err := json.MarshalWrite(
			cmd.OutOrStdout(), result,
			jsontext.WithIndent(" "), json.Deterministic(true),
		); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout()); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
	}

	return nil
}
