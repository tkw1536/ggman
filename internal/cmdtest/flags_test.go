//spellchecker:words cmdtest
package cmdtest_test

//spellchecker:words testing ggman internal cmdtest testutil
import (
	"testing"

	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/internal/cmdtest"
	"go.tkw01536.de/ggman/internal/testutil"
)

type withOverlap struct {
	POverlapsShort   bool `long:"p-overlaps" short:"P"`
	HereOverlapsLong bool `long:"here"       short:"q"`
	DoesNotOverlap   bool `long:"no-overlap" short:"z"`
}

func (withOverlap) Description() ggman.Description {
	panic("never called")
}

func (withOverlap) Run(context ggman.Context) error {
	panic("never called")
}

type noOverlap struct {
	DoesNotOverlap bool `long:"no-overlap" short:"z"`
}

func (noOverlap) Description() ggman.Description {
	panic("never called")
}

func (noOverlap) Run(context ggman.Context) error {
	panic("never called")
}

func TestAssertNoFlagOverlap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cmd         ggman.Command
		wantMessage string
		wantHelper  bool
	}{
		{
			name:        "command with overlap should fail",
			cmd:         &withOverlap{},
			wantMessage: "got FlagOverlap = [P here], but wanted []",
			wantHelper:  true,
		},
		{
			name:        "command without overlap should pass",
			cmd:         &noOverlap{},
			wantMessage: "",
			wantHelper:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var r testutil.RecordingT
			cmdtest.AssertNoFlagOverlap(&r, tt.cmd)

			if tt.wantMessage != r.Message {
				t.Errorf("cmdtest.AssertNoFlagOverlap() message = %q, want = %q", r.Message, tt.wantMessage)
			}

			if tt.wantHelper != r.HelperCalled {
				t.Errorf("cmdtest.AssertNoFlagOverlap() helper = %v, want = %v", r.HelperCalled, tt.wantHelper)
			}
		})
	}
}
