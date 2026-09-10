package cmd_test

//spellchecker:words path filepath testing ggman internal mockenv
import (
	"os"
	"path/filepath"
	"testing"

	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
)

func TestCommandCanon(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	tests := []struct {
		name    string
		workDir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"git@gitforge.example/user/repo",
			"",
			[]string{"canon", "git@gitforge.example/user/repo"},

			0,
			"git@gitforge.example:user/repo.git\n",
			"",
		},

		{
			"git@gitforge.example/user/repo ssh://%@^/$.git",
			"",
			[]string{"canon", "git@gitforge.example/user/repo", "ssh://%@^/$.git"},

			0,
			"ssh://user@gitforge.example/repo.git\n",
			"",
		},

		{
			"ssh://git@gitforge.example/hello/world",
			"",
			[]string{"canon", "ssh://git@gitforge.example/hello/world"},

			0,
			"git@gitforge.example:hello/world.git\n",
			"",
		},

		{
			"ssh://git@gitforge.example/hello/world ssh://%@^/$.git",
			"",
			[]string{"canon", "ssh://git@gitforge.example/hello/world", "ssh://%@^/$.git"},

			0,
			"ssh://hello@gitforge.example/world.git\n",
			"",
		},

		{
			"user@server.com/repo",
			"",
			[]string{"canon", "user@server.com/repo"},

			0,
			"git@server.com:user/repo.git\n",
			"",
		},

		{
			"user@server.com/repo ssh://%@^/$.git",
			"",
			[]string{"canon", "user@server.com/repo", "ssh://%@^/$.git"},

			0,
			"ssh://user@server.com/repo.git\n",
			"",
		},

		{
			"ssh://user@server.com:1234/repo.git",
			"",
			[]string{"canon", "ssh://user@server.com:1234/repo.git"},

			0,
			"git@server.com:user/repo.git\n",
			"",
		},

		{
			"ssh://user@server.com:1234/repo.git ssh://%@^/$.git",
			"",
			[]string{"canon", "ssh://user@server.com:1234/repo.git", "ssh://%@^/$.git"},

			0,
			"ssh://user@server.com/repo.git\n",
			"",
		},

		{
			"ssh://user@server.com:1234/repo.git $$",
			"",
			[]string{"canon", "ssh://user@server.com:1234/repo.git", "$$"},

			0,
			"ssh://user@server.com:1234/repo.git\n",
			"",
		},

		{
			"tree url strips forge reference",
			"",
			[]string{"canon", "https://gitforge.example/hello/world/tree/main", "git@^:$.git"},

			0,
			"git@gitforge.example:hello/world.git\n",
			"",
		},

		{
			"tree url with no-forge-split",
			"",
			[]string{"canon", "--no-forge-split", "https://gitforge.example/hello/world/tree/main", "git@^:$.git"},

			0,
			"git@gitforge.example:hello/world/tree/main.git\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, tt.workDir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}

func TestCommandCanon_ReplaceDomain(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	// Create and write out a CANFILE for purposes of this test.
	CANFILE := filepath.Join(t.TempDir(), "canfile")
	mock.SetCanfile(CANFILE)

	if err := os.WriteFile(CANFILE, []byte(`
# anything under example.com/namespace should be rewritten to a custom domain
^example.com/namespace git@!custom.com:$.git

# default to a normal pattern
git@^:$.git
	`), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		workDir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"normal repository",
			"",
			[]string{"canon", "git@server.com/user/repo"},

			0,
			"git@server.com:user/repo.git\n",
			"",
		},

		{
			"replaces domain for repo in namespace",
			"",
			[]string{"canon", "git@example.com/namespace/repo"},

			0,
			"git@custom.com:namespace/repo.git\n",
			"",
		},
		{
			"does not replace repo for other repo on same domain",
			"",
			[]string{"canon", "git@example.com/other/repo"},

			0,
			"git@example.com:other/repo.git\n",
			"",
		},
		{
			"does not replace repo for repo with similar prefix",
			"",
			[]string{"canon", "git@example.com/namespace_2/repo"},

			0,
			"git@example.com:namespace_2/repo.git\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, tt.workDir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
