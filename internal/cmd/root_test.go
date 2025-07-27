package cmd_test

//spellchecker:words testing ggman internal mockenv
import (
	"testing"

	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
)

// This test invokes the command with the help flag.
// It shouldn't fail.
func Test_main_docs(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)
	code, _, _ := mock.Run(t, cmd.NewCommand, "", "", "--help")
	if code != 0 {
		t.Fail()
	}
}
