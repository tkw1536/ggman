package cmd_test

//spellchecker:words testing ggman internal mockenv
import (
	"testing"

	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words workdir

func TestCommandLsr(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	mock.Clone("https://github.com/hello/world.git", "github.com", "hello", "world")
	mock.Clone("user@server.com/repo", "server.com", "user", "repo")
	mock.Clone("https://gitlab.com/hello/world.git", "gitlab.com", "hello", "world")

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"list remotes of all repositories",
			"",
			[]string{"lsr"},

			0,
			"https://github.com/hello/world.git\nhttps://gitlab.com/hello/world.git\nuser@server.com/repo\n",

			"",
		},

		{
			"list canonical remotes of all repositories",
			"",
			[]string{"lsr", "--canonical"},

			0,
			"git@github.com:hello/world.git\ngit@gitlab.com:hello/world.git\ngit@server.com:user/repo.git\n",

			"",
		},

		{
			"list remotes only hello/world repositories",
			"",
			[]string{"--for", "hello/world", "lsr"},

			0,
			"https://github.com/hello/world.git\nhttps://gitlab.com/hello/world.git\n",

			"",
		},

		{
			"list canonical remotes only hello/world repositories",
			"",
			[]string{"--for", "hello/world", "lsr", "--canonical"},

			0,
			"git@github.com:hello/world.git\ngit@gitlab.com:hello/world.git\n",

			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
