//spellchecker:words mockenv
package mockenv

//spellchecker:words testing essio shellescape github ggman internal testutil
import (
	"fmt"
	"testing"

	"al.essio.dev/pkg/shellescape"
	"github.com/tkw1536/ggman/internal/testutil"
)

// recordingT records a message passed to Errorf().
type recordingT struct {
	message string
}

func (f *recordingT) Errorf(format string, args ...any) {
	f.message = fmt.Sprintf(format, args...)
}

func TestMockEnv_AssertOutput(t *testing.T) {
	type fields struct {
		localRoot string
	}
	type args struct {
		prefix string
		got    string
		wants  []string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMessage string
	}{
		{"no replacement equal", fields{testutil.ToOSPath("/root/")}, args{"logprefix", "example", []string{"example"}}, ""},
		{"no replacement not equal", fields{testutil.ToOSPath("/root/")}, args{"logprefix", "example", []string{"example2"}}, "logprefix got = \"example\", want = \"example2\""},

		{"replace only ggroot ok", fields{testutil.ToOSPath("/root/")}, args{"logprefix", "prefix " + testutil.ToOSPath("/root") + " suffix", []string{"prefix ${GGROOT} suffix"}}, ""},
		{"replace only ggroot not ok", fields{testutil.ToOSPath("/root/")}, args{"logprefix", "prefix " + testutil.ToOSPath("/root") + " suffix", []string{"prefix ${GGROOT}/sub suffix"}}, fmt.Sprintf("logprefix got = %q, want = %q", "prefix "+testutil.ToOSPath("/root")+" suffix", "prefix "+testutil.ToOSPath("/root")+"/sub suffix")},

		{"replace full path ok", fields{testutil.ToOSPath("/root/")}, args{"logprefix", "prefix " + testutil.ToOSPath("/root/a/b") + " suffix", []string{"prefix ${GGROOT a b} suffix"}}, ""},
		{"replace full path not ok", fields{testutil.ToOSPath("/root/")}, args{"logprefix", "prefix " + testutil.ToOSPath("/root") + " suffix", []string{"prefix ${GGROOT a b} suffix"}}, fmt.Sprintf("logprefix got = %q, want = %q", "prefix "+testutil.ToOSPath("/root")+" suffix", "prefix "+testutil.ToOSPath("/root/a/b")+" suffix")},

		{"escape path with quotes", fields{testutil.ToOSPath("/root/")}, args{"logprefix", fmt.Sprintf("%q", testutil.ToOSPath("/root")), []string{"\"${GGROOT}\""}}, ""},
		{"escape path with `s", fields{testutil.ToOSPath("/!root/")}, args{"logprefix", shellescape.Quote(testutil.ToOSPath("/!root")), []string{"`${GGROOT}`"}}, ""},

		{"equal to first want is ok", fields{testutil.ToOSPath("/root/")}, args{"logprefix", "first", []string{"first", "last"}}, ""},
		{"equal to last want is ok", fields{testutil.ToOSPath("/root/")}, args{"logprefix", "last", []string{"first", "last"}}, ""},
		{"not equal to any wants shows last error", fields{testutil.ToOSPath("/root/")}, args{"logprefix", "neither", []string{"first error", "${GGROOT last}"}}, fmt.Sprintf("logprefix got = %q, want = %q", "neither", testutil.ToOSPath("/root/last"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockEnv{
				localRoot: tt.fields.localRoot,
			}

			var r recordingT
			mock.AssertOutput(&r, tt.args.prefix, tt.args.got, tt.args.wants...)

			if tt.wantMessage != r.message {
				t.Errorf("mock.AssertOutput() message = %q, want = %q", r.message, tt.wantMessage)
			}
		})
	}
}
