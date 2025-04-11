package cmd

//spellchecker:words testing github ggman internal mockenv
import (
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

//spellchecker:words shellrc

func TestCommandShellRC(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	code, stdout, stderr := mock.Run(Shellrc, "", "", "shellrc")
	if code != 0 {
		t.Errorf("Code = %d, wantCode = %d", code, 0)
	}
	if stdout != shellrcSh {
		t.Errorf("Got stdout = %s, expected = %s", stdout, shellrcSh)
	}
	if stderr != "" {
		t.Errorf("Got stderr = %s, expected = %s", stderr, "")
	}
}
