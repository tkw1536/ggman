package cmd_test

//spellchecker:words path filepath testing ggman internal mockenv
import (
	"os"
	"path/filepath"
	"testing"

	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words GGROOT workdir

func TestCommandFindFile(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	// with file 'example.txt'
	{
		clonePath := mock.Clone("https://github.com/hello/world.git", "github.com", "hello", "world")
		if err := os.WriteFile(filepath.Join(clonePath, "example.txt"), nil, 0600); err != nil {
			panic(err)
		}
	}

	// with file 'example/example.txt'
	{
		clonePath := mock.Clone("user@server.com/repo", "server.com", "user", "repo")
		if err := os.Mkdir(filepath.Join(clonePath, "example"), 0750); err != nil {
			panic(err)
		}
		if err := os.WriteFile(filepath.Join(clonePath, "example", "example.txt"), nil, 0600); err != nil {
			panic(err)
		}
	}

	// with nothing
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
			"find example.txt file",
			"",
			[]string{"find-file", "example.txt"},

			0,
			"${GGROOT github.com hello world}\n",
			"",
		},
		{
			"find example.txt file with paths",
			"",
			[]string{"find-file", "--print-file", "example.txt"},

			0,
			"${GGROOT github.com hello world example.txt}\n",
			"",
		},
		{
			"find example directory",
			"",
			[]string{"find-file", "example"},

			0,
			"${GGROOT server.com user repo}\n",
			"",
		},
		{
			"find example/example.txt file",
			"",
			[]string{"find-file", "example/example.txt"},

			0,
			"${GGROOT server.com user repo}\n",
			"",
		},
		{
			"don't find non-existent file",
			"",
			[]string{"find-file", "iDoNotExist.txt"},

			0,
			"",
			"",
		},
		{
			"don't find non-existent file with exit code",
			"",
			[]string{"find-file", "--exit-code", "iDoNotExist.txt"},

			1,
			"",
			"",
		},
		{
			"find existent file with exit code",
			"",
			[]string{"find-file", "--exit-code", "example.txt"},

			0,
			"${GGROOT github.com hello world}\n",
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
