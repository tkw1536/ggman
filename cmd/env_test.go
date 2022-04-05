package cmd

import (
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

func TestCommandRoot(t *testing.T) {
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
			"BUILDTIME\nGGROOT\nGIT\nPWD\nVERSION\n",
			"",
		},
		{
			"describe all variables",
			"",
			[]string{"env", "--describe"},

			0,
			"BUILDTIME: version of go this program was built with\nGGROOT: root folder all ggman repositories will be cloned to\nGIT: path to the native git\nPWD: current working directory\nVERSION: current ggman version\n",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Env, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
