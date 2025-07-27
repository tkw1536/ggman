package cmd_test

//spellchecker:words path filepath testing ggman internal mockenv
import (
	"os"
	"path/filepath"
	"testing"

	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words workdir GGROOT wrld

var testInputFile = `
; this and the following lines are ignored
# gitlab.com/hello/world
` + filepath.Join("server.com", "user", "repo") + `

// blank lines too

https://github.com/hello/world.git

`

func TestCommandLs(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	ghHelloWorld := mock.Clone("https://github.com/hello/world.git", "github.com", "hello", "world")
	serverRepo := mock.Clone("user@server.com/repo", "server.com", "user", "repo")
	glHelloWorld := mock.Clone("https://gitlab.com/hello/world.git", "gitlab.com", "hello", "world")

	inputFile := mock.Resolve("file.txt")
	if err := os.WriteFile(inputFile, []byte(testInputFile), 0600); err != nil {
		panic(err)
	}

	// make glHelloWorldDirty
	if err := os.WriteFile(filepath.Join(glHelloWorld, "dirty"), []byte{}, 0600); err != nil {
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
			"list all repositories with specific count",
			"",
			[]string{"ls", "--count", "2"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n",

			"",
		},

		{
			"list all repositories with higher than available count",
			"",
			[]string{"ls", "--count", "5"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",

			"",
		},

		{
			"don't support both one and count at the same time",
			"",
			[]string{"ls", "--one", "--count", "2"},

			4,
			"",

			"only one of `--one` and `--count` may be provided\n",
		},

		{
			"don't support negative limit",
			"",
			[]string{"ls", "--count", "-1"},

			4,
			"",

			"`--count` may not be negative\n",
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
			"list repositories with start flags",
			"",
			[]string{"--for", "^github.com", "ls"},

			0,
			"${GGROOT github.com hello world}\n",

			"",
		},

		{
			"list repositories with scores",
			"",
			[]string{"--for", "wrld", "ls", "--scores"},

			0,
			"0.900000 ${GGROOT github.com hello world}\n0.900000 ${GGROOT gitlab.com hello world}\n",

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
		{
			"list from input file",
			"",
			[]string{"--from-file", inputFile, "ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT server.com user repo}\n",

			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, stdout, stderr := mock.Run(t, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
