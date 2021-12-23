package cmd

import (
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

func TestCommandComps(t *testing.T) {
	mock := mockenv.NewMockEnv(t)

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"git@github.com/user/repo",
			"",
			[]string{"comps", "git@github.com/user/repo"},

			0,
			"github.com\nuser\nrepo\n",
			"",
		},

		{
			"ssh://git@github.com/hello/world",
			"",
			[]string{"comps", "ssh://git@github.com/hello/world"},

			0,
			"github.com\nhello\nworld\n",
			"",
		},

		{
			"user@server.com/repo",
			"",
			[]string{"comps", "user@server.com/repo"},

			0,
			"server.com\nuser\nrepo\n",
			"",
		},

		{
			"ssh://user@server.com:1234/repo.git",
			"",
			[]string{"comps", "ssh://user@server.com:1234/repo.git"},

			0,
			"server.com\nuser\nrepo\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Comps, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
