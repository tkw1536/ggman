package cmd_test

//spellchecker:words strconv testing github ggman internal mockenv
import (
	"strconv"
	"testing"

	"go.tkw01536.de/ggman/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words workdir GGROOT nolint tparallel paralleltest

//nolint:tparallel,paralleltest
func TestCommandLink(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	externalRepo := mock.Clone("https://github.com/hello/world.git", "..", "external")

	escapedExternalRepo := strconv.Quote(externalRepo)

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"Linking external repo",
			externalRepo,
			[]string{"link", "."},

			0,
			"Linking \"${GGROOT github.com hello world}\" -> " + escapedExternalRepo + "\n",
			"",
		},

		{
			"Linking external when it already exists",
			externalRepo,
			[]string{"link", "."},

			1,
			"",
			"another directory already exists in target location\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(cmd.Link, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
