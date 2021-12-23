package cmd

import (
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

func TestCommandLsr(t *testing.T) {
	mock := mockenv.NewMockEnv(t)

	mock.Register("https://github.com/hello/world.git")
	mock.Install("https://github.com/hello/world.git", "github.com", "hello", "world")

	mock.Register("user@server.com/repo")
	mock.Install("user@server.com/repo", "server.com", "user", "repo")

	mock.Register("https://gitlab.com/hello/world.git")
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
			code, stdout, stderr := mock.Run(Lsr, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
