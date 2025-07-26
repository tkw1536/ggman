//spellchecker:words main
package main

//spellchecker:words testing ggman internal mockenv
import (
	"testing"

	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words doccheck

// This test runs every command once with the --help flag.
// This should never fail, and have no effect.
func Test_main_docs(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	for _, name := range ggmanExe.Commands() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd, ok := ggmanExe.Command(name)
			if !ok {
				t.Fail()
			}
			code, _, _ := mock.RunLegacy(cmd, "", "", name, "--help")
			if code != 0 {
				t.Fail()
			}
		})
	}
}
