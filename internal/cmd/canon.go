package cmd

//spellchecker:words github cobra ggman internal pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words github cobra ggman goprogram exit cspec

func NewCanonCommand() *cobra.Command {
	impl := new(canon)

	cmd := &cobra.Command{
		Use:   "canon URL [CANSPEC]",
		Short: "Print the canonical version of a URL",
		Long: `Canon prints the canonical version of a URL.
An optional second argument specifies the CANSPEC.

On github.com and other forges, repositories can be cloned via multiple URLs:

- https://github.com/hello/world.git
- git@github.com:hello/world.git

The SSH URL is typically preferred to avoid password entry.
ggman treats the canonical URL as primary and uses it for cloning.

A CANSPEC (canonical specification) transforms URLs into canonical form.
An example CANSPEC is 'git@^:$.git'.

CANSPECs work by splitting URLs into path-like components with normalization:

- 'git@github.com/user/repo' => 'github.com', 'user', 'repo'
- 'github.com/hello/world.git' => 'github.com', 'hello', 'world'
- 'user@server.com:repo.git' => 'server.com', 'user', 'repo'

The 'ggman comps' command shows URL components.

CANSPEC characters are copied literally except:

- '^' is replaced by the first component (hostname)
- '%' is replaced by the second unused component (commonly username)
- '!' consumes a component without using it in the output
- '$' is replaced by remaining components joined with '/'; stops special processing

If '$' is absent, it is assumed at the end.
The CANSPEC '$$' leaves URLs unchanged.

Examples for components 'server.com', 'user', 'repository':

- 'git@^:$.git' => 'git@server.com:user/repository.git'
- 'ssh://%@^/$.git' => 'ssh://user@server.com/repository.git'
- 'git@!ssh.server.com:$.git' => 'git@ssh.server.com:user/repository.git'
- empty CANSPEC => 'server.com/user/repository'

A CANFILE provides pattern-based CANSPEC rules.
It is loaded from '.ggman' in the home directory or from '$GGMAN_CANFILE'.

Each CANFILE line contains one or two space-separated strings: a pattern and an optional CANSPEC.
Lines starting with '#' or '\' are comments.

Example CANFILE:

    # for anything on git.example.com, clone with https
    ^git.example.com https://$.git

    # for anything under a specific namespace use a custom domain name
    ^git2.example.com/my_namespace git@!ssh.example.com:$.git

    # for anything else on git2.example.com leave the urls unchanged
    ^git2.example.com $$

    # by default, clone via ssh
    git@^:$.git

Omitting the CANSPEC argument uses the CANFILE for resolution.`,
		Args: cobra.RangeArgs(1, 2),

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type canon struct {
	Positional struct {
		URL     env.URL
		CANSPEC string
	}
}

var (
	errCanonUnableCanFile = exit.NewErrorWithCode("unable to load default CANFILE", env.ExitContext)
)

func (c *canon) ParseArgs(cmd *cobra.Command, args []string) error {
	c.Positional.URL = env.ParseURL(args[0])
	if len(args) == 2 {
		c.Positional.CANSPEC = args[1]
	}
	return nil
}

func (c *canon) Exec(cmd *cobra.Command, args []string) error {
	var file env.CanFile

	if c.Positional.CANSPEC == "" {
		env, err := env.GetEnv(cmd, env.Requirement{})
		if err != nil {
			return fmt.Errorf("failed to get environment: %w", err)
		}

		if file, err = env.LoadDefaultCANFILE(); err != nil {
			return fmt.Errorf("%w: %w", errCanonUnableCanFile, err)
		}
	} else {
		file = []env.CanLine{{Pattern: "", Canonical: c.Positional.CANSPEC}}
	}

	// print out the canonical version of the file
	canonical := c.Positional.URL.CanonicalWith(file)

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), canonical); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	return nil
}
