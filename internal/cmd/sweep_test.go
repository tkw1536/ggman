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

func TestCommandSweep(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	path := mock.Clone(t.Context(), "https://gitforge.example/hello/world.git", "gitforge.example", "hello", "world")
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
	mkdir(filepath.Join("gitforge.example", "hello", "world", "empty"))
	mkdir(filepath.Join("gitforge.example", "empty", "empty1"))
	mkdir(filepath.Join("gitforge.example", "empty", "empty2"))
	mkdir(filepath.Join("gitforge.example", "full"), "file")

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
			"${GGROOT gitforge.example empty empty1}\n${GGROOT gitforge.example empty empty2}\n${GGROOT gitforge.example empty}\n",
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
