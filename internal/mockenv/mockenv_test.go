// Package mockenv contains facilities for unit testing commands
//
//spellchecker:words mockenv
package mockenv_test

//spellchecker:words path filepath testing ggman internal mockenv pkglib testlib
import (
	"testing"

	"go.tkw01536.de/ggman/internal/mockenv"
	"go.tkw01536.de/pkglib/testlib"
)

//spellchecker:words GGROOT logprefix

func TestMockEnv_Register(t *testing.T) {
	t.Parallel()

	const remote = "https://examaple.com/repo.git"

	mock := mockenv.NewMockEnv(t)
	mock.Register(remote)

	panicked, _ := testlib.DoesPanic(func() {
		mock.Register(remote)
	})

	if !panicked {
		t.Errorf("MockEnv.Register: Allowed dual registration")
	}
}
