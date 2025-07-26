package cmd_test

//spellchecker:words testing ggman internal mockenv
import (
	"testing"

	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words workdir GGROOT

func TestCommandEnv(t *testing.T) {
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
			"list all variables",
			"",
			[]string{"env", "--list"},

			0,
			"GGMAN_TIME\nGGMAN_VERSION\nGGROOT\nGIT\nPWD\n",
			"",
		},
		{
			"describe all variables",
			"",
			[]string{"env", "--describe"},

			0,
			"GGMAN_TIME: the time this version of ggman was built\nGGMAN_VERSION: the version of ggman this version is\nGGROOT: root folder all ggman repositories will be cloned to\nGIT: path to the native git\nPWD: current working directory\n",
			"",
		},

		{
			"show single variable",
			"",
			[]string{"env", "GGROOT"},

			0,
			"GGROOT=`${GGROOT}`\n",
			"",
		},
		{
			"show unknown variable",
			"",
			[]string{"env", "NOT_EXISTS"},

			4,
			"",
			"\"NOT_EXISTS\": unknown environment variable\n",
		},
		{
			"list single variable",
			"",
			[]string{"env", "--list", "GGROOT"},

			0,
			"GGROOT\n",
			"",
		},
		{
			"list single variable (case-insensitive)",
			"",
			[]string{"env", "--list", "ggroot"},

			0,
			"GGROOT\n",
			"",
		},
		{
			"describe single variable",
			"",
			[]string{"env", "--describe", "GGROOT"},

			0,
			"GGROOT: root folder all ggman repositories will be cloned to\n",
			"",
		},
		{
			"raw single variable",
			"",
			[]string{"env", "--raw", "GGROOT"},

			0,
			"${GGROOT}\n",
			"",
		},
		{
			"more than one mode",
			"",
			[]string{"env", "--list", "--raw"},

			4,
			"",
			"at most one of `--raw`, `--list` and `--describe` may be given\n",
		},
		{
			"more than one mode (2)",
			"",
			[]string{"env", "--list", "--describe"},

			4,
			"",
			"at most one of `--raw`, `--list` and `--describe` may be given\n",
		},
		{
			"more than one mode (3)",
			"",
			[]string{"env", "--raw", "--describe"},

			4,
			"",
			"at most one of `--raw`, `--list` and `--describe` may be given\n",
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
