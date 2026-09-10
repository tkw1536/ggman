package cmd_test

//spellchecker:words testing ggman internal mockenv
import (
	"testing"

	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words ggman GGROOT workdir

func TestCommandWhere(t *testing.T) {
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
			"show directory of repository",
			"",
			[]string{"where", "https://gitforge.example/hello/world.git"},

			0,
			"${GGROOT gitforge.example hello world}\n",
			"",
		},

		{
			"tree url strips forge reference",
			"",
			[]string{"where", "https://gitforge.example/hello/world/tree/main"},

			0,
			"${GGROOT gitforge.example hello world}\n",
			"",
		},

		{
			"tree url with no-forge-split",
			"",
			[]string{"where", "--no-forge-split", "https://gitforge.example/hello/world/tree/main"},

			0,
			"${GGROOT gitforge.example hello world tree main}\n",
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
