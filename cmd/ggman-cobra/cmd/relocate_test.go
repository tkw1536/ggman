package cmd_test

//spellchecker:words path filepath testing ggman internal mockenv
import (
	"os"
	"path/filepath"
	"testing"

	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words workdir GGROOT nolint tparallel paralleltest

//nolint:tparallel,paralleltest
func TestCommandRelocate(t *testing.T) {
	t.Parallel()

	symlink := func(oldName, newName string) {
		err := os.Symlink(oldName, newName)
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

	// link in an external repository in the right place
	external2 := mock.Clone("https://github.com/right/external2.git", "..", "external-path-2")
	symlink(external2, mock.Resolve(filepath.Join("github.com", "right", "wrong-external")))

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
			"mkdir -p `${GGROOT github.com right}`\nmv `${GGROOT github.com right wrong-external}` `${GGROOT github.com right external2}`\nmkdir -p `${GGROOT github.com correct}`\nmv `${GGROOT github.com incorrect directory}` `${GGROOT github.com correct directory}`\n",

			"",
		},

		{
			"relocate without simulate",
			"",
			[]string{"relocate"},

			0,
			"mkdir -p `${GGROOT github.com right}`\nmv `${GGROOT github.com right wrong-external}` `${GGROOT github.com right external2}`\nmkdir -p `${GGROOT github.com correct}`\nmv `${GGROOT github.com incorrect directory}` `${GGROOT github.com correct directory}`\n",

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
			code, stdout, stderr := mock.Run(t, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}

func TestCommandRelocate_existsRepo(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	// clone the same repository twice
	mock.Register("https://github.com/right/directory.git")
	mock.Install("https://github.com/right/directory.git", "github.com", "right", "directory")
	mock.Install("https://github.com/right/directory.git", "github.com", "right", "other")

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
			"mkdir -p `${GGROOT github.com right}`\nmv `${GGROOT github.com right other}` `${GGROOT github.com right directory}`\n",

			"",
		},

		{
			"relocate without simulate",
			"",
			[]string{"relocate"},

			1,
			"mkdir -p `${GGROOT github.com right}`\nmv `${GGROOT github.com right other}` `${GGROOT github.com right directory}`\n",

			"repository already exists at \"${GGROOT github.com right directory}\"\n",
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

func TestCommandRelocate_existsPath(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	// clone the same repository twice
	mock.Clone("https://github.com/right/directory.git", "github.com", "wrong", "directory")

	if err := os.MkdirAll(mock.Resolve("github.com", "right", "directory"), os.ModePerm|os.ModeDir); err != nil {
		panic(err)
	}

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
			"mkdir -p `${GGROOT github.com right}`\nmv `${GGROOT github.com wrong directory}` `${GGROOT github.com right directory}`\n",

			"",
		},

		{
			"relocate without simulate",
			"",
			[]string{"relocate"},

			1,
			"mkdir -p `${GGROOT github.com right}`\nmv `${GGROOT github.com wrong directory}` `${GGROOT github.com right directory}`\n",

			"\"${GGROOT github.com right directory}\": path already exists\n",
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
