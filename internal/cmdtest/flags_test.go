package cmdtest_test

import (
	"testing"

	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/internal/cmdtest"
	"go.tkw01536.de/ggman/internal/testutil"
)

type withOverlap struct {
	POverlapsShort   bool `long:"p-overlaps" short:"p"`
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

func TestAssertFlagOverlap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cmd  ggman.Command
		want []string

		wantMessage string
		wantHelper  bool
	}{
		{
			name: "wrong check on command with overlap",
			cmd:  &withOverlap{},
			want: []string{},

			wantMessage: "got FlagOverlap = [here p], but wanted []",
			wantHelper:  true,
		},

		{
			name: "correct check on command with overlap",
			cmd:  &withOverlap{},
			want: []string{"here", "p"},

			wantMessage: "",
			wantHelper:  true,
		},

		{
			name: "wrong check on command without overlap",
			cmd:  &noOverlap{},
			want: []string{"a", "b"},

			wantMessage: "got FlagOverlap = [], but wanted [a b]",
			wantHelper:  true,
		},

		{
			name: "correct check on command without overlap",
			cmd:  &noOverlap{},
			want: []string{},

			wantMessage: "",
			wantHelper:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var r testutil.RecordingT
			cmdtest.AssertFlagOverlap(&r, tt.cmd, tt.want)

			if tt.wantMessage != r.Message {
				t.Errorf("cmdtest.AssertFlagOverlap() message = %q, want = %q", r.Message, tt.wantMessage)
			}

			if tt.wantHelper != r.HelperCalled {
				t.Errorf("cmdtest.AssertFlagOverlap() helper = %v, want = %v", r.HelperCalled, tt.wantHelper)
			}
		})
	}
}
