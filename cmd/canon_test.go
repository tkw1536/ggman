package cmd

//spellchecker:words testing github ggman internal mockenv
import (
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

func TestCommandCanon(t *testing.T) {
	mock := mockenv.NewMockEnv(t)

	tests := []struct {
		name    string
		workDir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"git@github.com/user/repo",
			"",
			[]string{"canon", "git@github.com/user/repo"},

			0,
			"git@github.com:user/repo.git\n",
			"",
		},

		{
			"git@github.com/user/repo ssh://%@^/$.git",
			"",
			[]string{"canon", "git@github.com/user/repo", "ssh://%@^/$.git"},

			0,
			"ssh://user@github.com/repo.git\n",
			"",
		},

		{
			"ssh://git@github.com/hello/world",
			"",
			[]string{"canon", "ssh://git@github.com/hello/world"},

			0,
			"git@github.com:hello/world.git\n",
			"",
		},

		{
			"ssh://git@github.com/hello/world ssh://%@^/$.git",
			"",
			[]string{"canon", "ssh://git@github.com/hello/world", "ssh://%@^/$.git"},

			0,
			"ssh://hello@github.com/world.git\n",
			"",
		},

		{
			"user@server.com/repo",
			"",
			[]string{"canon", "user@server.com/repo"},

			0,
			"git@server.com:user/repo.git\n",
			"",
		},

		{
			"user@server.com/repo ssh://%@^/$.git",
			"",
			[]string{"canon", "user@server.com/repo", "ssh://%@^/$.git"},

			0,
			"ssh://user@server.com/repo.git\n",
			"",
		},

		{
			"ssh://user@server.com:1234/repo.git",
			"",
			[]string{"canon", "ssh://user@server.com:1234/repo.git"},

			0,
			"git@server.com:user/repo.git\n",
			"",
		},

		{
			"ssh://user@server.com:1234/repo.git ssh://%@^/$.git",
			"",
			[]string{"canon", "ssh://user@server.com:1234/repo.git", "ssh://%@^/$.git"},

			0,
			"ssh://user@server.com/repo.git\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Canon, tt.workDir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
