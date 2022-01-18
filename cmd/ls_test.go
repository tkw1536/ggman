package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

func TestCommandLs(t *testing.T) {
	mock := mockenv.NewMockEnv(t)

	ghHelloWorld := mock.Clone("https://github.com/hello/world.git", "github.com", "hello", "world")
	serverRepo := mock.Clone("user@server.com/repo", "server.com", "user", "repo")
	glHelloWorld := mock.Clone("https://gitlab.com/hello/world.git", "gitlab.com", "hello", "world")

	// make glHelloWorldDirty
	if err := os.WriteFile(filepath.Join(glHelloWorld, "dirty"), []byte{}, os.ModePerm); err != nil {
		panic(err)
	}

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
			"list dirty and clean repositories",
			"",
			[]string{"--dirty", "--clean", "ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",

			"",
		},

		{
			"list dirty repositories only",
			"",
			[]string{"--dirty", "ls"},

			0,
			"${GGROOT gitlab.com hello world}\n",

			"",
		},

		{
			"list clean repositories only",
			"",
			[]string{"--clean", "ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT server.com user repo}\n",

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
			"list only clean hello/world repositories",
			"",
			[]string{"--for", "hello/world", "--clean", "ls"},

			0,
			"${GGROOT github.com hello world}\n",

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
			"list an absolute path with --path",
			serverRepo,
			[]string{"--path", ghHelloWorld, "ls"},

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

		{
			"list a relative path with --path",
			glHelloDir,
			[]string{"--path", filepath.Join(".", "world"), "ls"},

			0,
			"${GGROOT gitlab.com hello world}\n",

			"",
		},

		{
			"list multiple paths with --path",
			glHelloDir,
			[]string{"--path", filepath.Join(".", "world"), "--path", ghHelloWorld, "ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n",

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
