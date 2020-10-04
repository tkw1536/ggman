package cmd

import (
	"testing"

	"github.com/tkw1536/ggman/testutil/mockenv"
)

func TestCommandLs(t *testing.T) {
	mock, cleanup := mockenv.NewMockEnv()
	defer cleanup()

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
			"list all repositories",
			"",
			[]string{"ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",

			"",
		},

		{
			"list all repositories with exit code",
			"",
			[]string{"ls", "--exit-code"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",

			"",
		},

		{
			"list only hello/world repositories",
			"",
			[]string{"--for", "hello/world", "ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n",

			"",
		},

		{
			"list non-existing repositories",
			"",
			[]string{"--for", "does/not/exist", "ls"},

			0,
			"",

			"",
		},

		{
			"list non-existing repositories with exit code",
			"",
			[]string{"--for", "does/not/exist", "ls", "--exit-code"},

			1,
			"",

			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Ls, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, stdout, tt.wantStdout, "Stdout")
			mock.AssertOutput(t, stderr, tt.wantStderr, "Stderr")
		})
	}
}
