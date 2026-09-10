package cmd_test

//spellchecker:words testing ggman internal mockenv
import (
	"testing"

	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words workdir

func TestCommandParse(t *testing.T) {
	t.Parallel()

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
			"comps: git@gitforge.example/user/repo",
			"",
			[]string{"parse", "--comps", "git@gitforge.example/user/repo"},

			0,
			"gitforge.example\nuser\nrepo\n",
			"",
		},

		{
			"comps: ssh://git@gitforge.example/hello/world",
			"",
			[]string{"parse", "--comps", "ssh://git@gitforge.example/hello/world"},

			0,
			"gitforge.example\nhello\nworld\n",
			"",
		},

		{
			"comps: user@server.com/repo",
			"",
			[]string{"parse", "--comps", "user@server.com/repo"},

			0,
			"server.com\nuser\nrepo\n",
			"",
		},

		{
			"comps: ssh://user@server.com:1234/repo.git",
			"",
			[]string{"parse", "--comps", "ssh://user@server.com:1234/repo.git"},

			0,
			"server.com\nuser\nrepo\n",
			"",
		},

		{
			"comps: tree url strips forge reference",
			"",
			[]string{"parse", "--comps", "https://gitforge.example/hello/world/tree/main/src"},

			0,
			"gitforge.example\nhello\nworld\n",
			"",
		},

		{
			"json default with tree url",
			"",
			[]string{"parse", "https://gitforge.example/hello/world/tree/main/src"},

			0,
			"{\n \"comps\": [\n  \"gitforge.example\",\n  \"hello\",\n  \"world\"\n ],\n \"branch\": \"main\",\n \"path\": \"src\"\n}\n",
			"",
		},

		{
			"json without tree",
			"",
			[]string{"parse", "https://gitforge.example/hello/world.git"},

			0,
			"{\n \"comps\": [\n  \"gitforge.example\",\n  \"hello\",\n  \"world\"\n ],\n \"branch\": \"\",\n \"path\": \"\"\n}\n",
			"",
		},

		{
			"branch only",
			"",
			[]string{"parse", "--branch", "https://gitforge.example/hello/world/tree/main/src"},

			0,
			"main\n",
			"",
		},

		{
			"path only",
			"",
			[]string{"parse", "--path", "https://gitforge.example/hello/world/tree/main/src"},

			0,
			"src\n",
			"",
		},

		{
			"no-forge-split keeps tree in comps",
			"",
			[]string{"parse", "--comps", "--no-forge-split", "https://gitforge.example/hello/world/tree/main/src"},

			0,
			"gitforge.example\nhello\nworld\ntree\nmain\nsrc\n",
			"",
		},

		{
			"no-forge-split empties branch",
			"",
			[]string{"parse", "--branch", "--no-forge-split", "https://gitforge.example/hello/world/tree/main/src"},

			0,
			"\n",
			"",
		},

		{
			"exclusive comps and branch",
			"",
			[]string{"parse", "--comps", "--branch", "https://gitforge.example/hello/world.git"},

			4,
			"",
			"only one of \"--comps\", \"--branch\" and \"--path\" may be provided\n",
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

func TestCommandCompsAlias(t *testing.T) {
	t.Parallel()

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
			"git@gitforge.example/user/repo",
			"",
			[]string{"comps", "git@gitforge.example/user/repo"},

			0,
			"gitforge.example\nuser\nrepo\n",
			"",
		},

		{
			"ssh://git@gitforge.example/hello/world",
			"",
			[]string{"comps", "ssh://git@gitforge.example/hello/world"},

			0,
			"gitforge.example\nhello\nworld\n",
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
