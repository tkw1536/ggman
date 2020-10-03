package cmd

import (
	"testing"

	"github.com/tkw1536/ggman/testutil/mockenv"
)

func TestCommandLink(t *testing.T) {
	mock, cleanup := mockenv.NewMockEnv()
	defer cleanup()

	mock.Register("https://github.com/hello/world.git")
	externalRepo := mock.Install("https://github.com/hello/world.git", "..", "external")

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
			"Linking \"${GGROOT github.com hello world}\" -> \"" + externalRepo + "\"\n",
			"",
		},

		{
			"Linking external when it already exists",
			externalRepo,
			[]string{"link", "."},

			1,
			"",
			"Unable to link repository: Another directory already exists in target location. \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Link, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, stdout, tt.wantStdout, "Stdout")
			mock.AssertOutput(t, stderr, tt.wantStderr, "Stderr")
		})
	}
}
