package cmd_test

//spellchecker:words testing ggman internal cmdtest
import (
	"testing"

	"go.tkw01536.de/ggman/cmd"
	"go.tkw01536.de/ggman/internal/cmdtest"
)

func TestCommandShellRc_Overlap(t *testing.T) {
	t.Parallel()

	cmdtest.AssertFlagOverlap(t, cmd.Shellrc, []string{})
}
