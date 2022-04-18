package cmd

import (
	"fmt"
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

func TestCommandLink(t *testing.T) {
	mock := mockenv.NewMockEnv(t)

	externalRepo := mock.Clone("https://github.com/hello/world.git", "..", "external")

	escapedExternalRepo := fmt.Sprintf("%q", externalRepo)

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
			"unable to link repository: another directory already exists in target location\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Link, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
