package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

func TestCommandRelocate(t *testing.T) {
	symlink := func(oldname, newname string) {
		err := os.Symlink(oldname, newname)
		if err != nil {
			panic(err)
		}
	}

	mock := mockenv.NewMockEnv(t)

	mock.Clone("https://github.com/right/directory.git", "github.com", "right", "directory")
	mock.Clone("https://github.com/correct/directory.git", "github.com", "incorrect", "directory")

	// link in an external repository in the right place
	external1 := mock.Clone("https://github.com/right/external1.git", "..", "external-path-1")
	symlink(external1, mock.Resolve(filepath.Join("github.com", "right", "external1")))

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"relocate with simulate",
			"",
			[]string{"relocate", "--simulate"},

			0,
			"mkdir -p `${GGROOT github.com correct}`\nmv `${GGROOT github.com incorrect directory}` `${GGROOT github.com correct directory}`\n",

			"",
		},

		{
			"relocate without simulate",
			"",
			[]string{"relocate"},

			0,
			"mkdir -p `${GGROOT github.com correct}`\nmv `${GGROOT github.com incorrect directory}` `${GGROOT github.com correct directory}`\n",

			"",
		},

		{
			"nothing to relocate",
			"",
			[]string{"relocate"},

			0,
			"",

			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Relocate, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
