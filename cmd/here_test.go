package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

func TestCommandHere(t *testing.T) {
	mock := mockenv.NewMockEnv(t)

	clonePath := mock.Clone("https://github.com/hello/world.git", "github.com", "hello", "world")

	subClonePath := filepath.Join(clonePath, "sub")
	os.MkdirAll(subClonePath, os.ModePerm)

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"Open url at root",
			clonePath,
			[]string{"here"},
			0,
			"${GGROOT github.com hello world}\n",
			"",
		},

		{
			"Print path with tree at root",
			clonePath,
			[]string{"here", "--tree"},
			0,
			"${GGROOT github.com hello world}\n.\n",
			"",
		},

		{
			"Print path at subpath",
			subClonePath,
			[]string{"here"},
			0,
			"${GGROOT github.com hello world}\n",
			"",
		},

		{
			"Print path with tree at subpath",
			subClonePath,
			[]string{"here", "--tree"},
			0,
			"${GGROOT github.com hello world}\nsub\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Here, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
