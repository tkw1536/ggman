package cmd_test

//spellchecker:words testing ggman internal mockenv
import (
	"testing"

	"go.tkw01536.de/ggman/cmd/ggman-cobra/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
)

func TestCommandShellRC(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	code, stdout, stderr := mock.Run(t, "", "", "shellrc")
	if code != 0 {
		t.Errorf("Code = %d, wantCode = %d", code, 0)
	}
	if stdout != cmd.ShellrcSh {
		t.Errorf("Got stdout = %s, expected = %s", stdout, cmd.ShellrcSh)
	}
	if stderr != "" {
		t.Errorf("Got stderr = %s, expected = %s", stderr, "")
	}
}
