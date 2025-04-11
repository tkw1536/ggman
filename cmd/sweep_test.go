package cmd_test

//spellchecker:words path filepath testing github ggman internal mockenv
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/ggman/cmd"
	"github.com/tkw1536/ggman/internal/mockenv"
)

//spellchecker:words GGROOT workdir

func TestCommandSweep(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	path := mock.Clone("https://github.com/hello/world.git", "github.com", "hello", "world")
	base := filepath.Join(path, "..", "..", "..")

	mkdir := func(s string, files ...string) {
		path := filepath.Join(base, s)
		err := os.MkdirAll(path, 0750)
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			if err := os.WriteFile(filepath.Join(path, f), nil, 0600); err != nil {
				panic(err)
			}
		}
	}
	mkdir(filepath.Join("github.com", "hello", "world", "empty"))
	mkdir(filepath.Join("github.com", "empty", "empty1"))
	mkdir(filepath.Join("github.com", "empty", "empty2"))
	mkdir(filepath.Join("github.com", "full"), "file")

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"list empty directories",
			"",
			[]string{"sweep"},

			0,
			"${GGROOT github.com empty empty1}\n${GGROOT github.com empty empty2}\n${GGROOT github.com empty}\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, stdout, stderr := mock.Run(cmd.Sweep, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
