package cmd_test

//spellchecker:words testing ggman internal cmdtest mockenv
import (
	"testing"

	"go.tkw01536.de/ggman/cmd"
	"go.tkw01536.de/ggman/internal/cmdtest"
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
			[]string{"where", "https://github.com/hello/world.git"},

			0,
			"${GGROOT github.com hello world}\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, stdout, stderr := mock.Run(cmd.Where, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}

func TestCommandWhere_Overlap(t *testing.T) {
	t.Parallel()

	cmdtest.AssertNoFlagOverlap(t, cmd.Where)
}
