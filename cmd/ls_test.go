package cmd

import (
	"path/filepath"
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

func TestCommandLs(t *testing.T) {
	mock, cleanup := mockenv.NewMockEnv()
	defer cleanup()

	mock.Register("https://github.com/hello/world.git")
	ghHelloWorld := mock.Install("https://github.com/hello/world.git", "github.com", "hello", "world")

	mock.Register("user@server.com/repo")
	serverRepo := mock.Install("user@server.com/repo", "server.com", "user", "repo")

	mock.Register("https://gitlab.com/hello/world.git")
	glHelloWorld := mock.Install("https://gitlab.com/hello/world.git", "gitlab.com", "hello", "world")

	glHelloDir := filepath.Join(glHelloWorld, "..")

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
			"list all repositories with one",
			"",
			[]string{"ls", "--one"},

			0,
			"${GGROOT github.com hello world}\n",

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
			"list repositories fuzzy",
			"",
			[]string{"--for", "wrld", "ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n",

			"",
		},

		{
			"list repositories non-fuzzy",
			"",
			[]string{"--no-fuzzy-filter", "--for", "wrld", "ls"},

			0,
			"",

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

		{
			"list only current repository (github.com hello world)",
			ghHelloWorld,
			[]string{"--here", "ls"},

			0,
			"${GGROOT github.com hello world}\n",

			"",
		},

		{
			"list only current repository (server.com user repo)",
			serverRepo,
			[]string{"--here", "ls"},

			0,
			"${GGROOT server.com user repo}\n",

			"",
		},
		{
			"list only current repository (gitlab.com hello world)",
			glHelloWorld,
			[]string{"--here", "ls"},

			0,
			"${GGROOT gitlab.com hello world}\n",

			"",
		},
		{
			"list an absolute path",
			serverRepo,
			[]string{"--for", ghHelloWorld, "ls"},

			0,
			"${GGROOT github.com hello world}\n",

			"",
		},
		{
			"list a relative path",
			glHelloDir,
			[]string{"--for", filepath.Join(".", "world"), "ls"},

			0,
			"${GGROOT gitlab.com hello world}\n",

			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Ls, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
