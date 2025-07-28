package cmd_test

//spellchecker:words testing ggman internal mockenv
import (
	"testing"

	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words GGROOT workdir tparallel paralleltest

//nolint:tparallel,paralleltest
func TestCommandFix(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	mock.Register("https://github.com/hello/world.git", "git@github.com:hello/world.git")
	mock.Install("https://github.com/hello/world.git", "github.com", "hello", "world")

	mock.Register("user@server.com/repo", "git@server.com:user/repo.git")
	mock.Install("user@server.com/repo", "server.com", "user", "repo")

	mock.Register("https://gitlab.com/hello/world.git", "git@gitlab.com:hello/world.git")
	mock.Install("https://gitlab.com/hello/world.git", "gitlab.com", "hello", "world")

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"simulate fixing remotes of all repositories",
			"",
			[]string{"fix", "--simulate"},

			0,
			"Simulate fixing remote of \"${GGROOT github.com hello world}\"Updating origin: https://github.com/hello/world.git -> git@github.com:hello/world.git\nSimulate fixing remote of \"${GGROOT gitlab.com hello world}\"Updating origin: https://gitlab.com/hello/world.git -> git@gitlab.com:hello/world.git\nSimulate fixing remote of \"${GGROOT server.com user repo}\"Updating origin: user@server.com/repo -> git@server.com:user/repo.git\n",
			"",
		},

		{
			"actually fixing remotes of all repositories",
			"",
			[]string{"fix"},

			0,
			"Fixing remote of \"${GGROOT github.com hello world}\"Updating origin: https://github.com/hello/world.git -> git@github.com:hello/world.git\nFixing remote of \"${GGROOT gitlab.com hello world}\"Updating origin: https://gitlab.com/hello/world.git -> git@gitlab.com:hello/world.git\nFixing remote of \"${GGROOT server.com user repo}\"Updating origin: user@server.com/repo -> git@server.com:user/repo.git\n",
			"",
		},

		{
			"fixing remotes of fixed repositories",
			"",
			[]string{"fix"},

			0,
			"",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
