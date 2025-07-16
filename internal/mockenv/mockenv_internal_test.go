//spellchecker:words mockenv
package mockenv

//spellchecker:words testing essio shellescape ggman internal testutil
import (
	"fmt"
	"testing"

	"al.essio.dev/pkg/shellescape"
	"go.tkw01536.de/ggman/internal/testutil"
)

//spellchecker:words logprefix ggroot GGROOT

// recordingT records a message passed to Errorf().
type recordingT struct {
	message string
	helper  bool
}

func (f *recordingT) Helper() {
	f.helper = true
}

func (f *recordingT) Errorf(format string, args ...any) {
	f.message = fmt.Sprintf(format, args...)
}

func TestMockEnv_AssertOutput(t *testing.T) {
	t.Parallel()

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
		wantHelper  bool
	}{
		{name: "no replacement equal", fields: fields{localRoot: testutil.ToOSPath("/root/")}, args: args{prefix: "logprefix", got: "example", wants: []string{"example"}}, wantMessage: "", wantHelper: true},
		{name: "no replacement not equal", fields: fields{localRoot: testutil.ToOSPath("/root/")}, args: args{prefix: "logprefix", got: "example", wants: []string{"example2"}}, wantMessage: "logprefix got = \"example\", want = \"example2\"", wantHelper: true},

		{name: "replace only ggroot ok", fields: fields{localRoot: testutil.ToOSPath("/root/")}, args: args{prefix: "logprefix", got: "prefix " + testutil.ToOSPath("/root") + " suffix", wants: []string{"prefix ${GGROOT} suffix"}}, wantMessage: "", wantHelper: true},
		{name: "replace only ggroot not ok", fields: fields{localRoot: testutil.ToOSPath("/root/")}, args: args{prefix: "logprefix", got: "prefix " + testutil.ToOSPath("/root") + " suffix", wants: []string{"prefix ${GGROOT}/sub suffix"}}, wantMessage: fmt.Sprintf("logprefix got = %q, want = %q", "prefix "+testutil.ToOSPath("/root")+" suffix", "prefix "+testutil.ToOSPath("/root")+"/sub suffix"), wantHelper: true},

		{name: "replace full path ok", fields: fields{localRoot: testutil.ToOSPath("/root/")}, args: args{prefix: "logprefix", got: "prefix " + testutil.ToOSPath("/root/a/b") + " suffix", wants: []string{"prefix ${GGROOT a b} suffix"}}, wantMessage: "", wantHelper: true},
		{name: "replace full path not ok", fields: fields{localRoot: testutil.ToOSPath("/root/")}, args: args{prefix: "logprefix", got: "prefix " + testutil.ToOSPath("/root") + " suffix", wants: []string{"prefix ${GGROOT a b} suffix"}}, wantMessage: fmt.Sprintf("logprefix got = %q, want = %q", "prefix "+testutil.ToOSPath("/root")+" suffix", "prefix "+testutil.ToOSPath("/root/a/b")+" suffix"), wantHelper: true},

		{name: "escape path with quotes", fields: fields{localRoot: testutil.ToOSPath("/root/")}, args: args{prefix: "logprefix", got: fmt.Sprintf("%q", testutil.ToOSPath("/root")), wants: []string{"\"${GGROOT}\""}}, wantMessage: "", wantHelper: true},
		{name: "escape path with `s", fields: fields{localRoot: testutil.ToOSPath("/!root/")}, args: args{prefix: "logprefix", got: shellescape.Quote(testutil.ToOSPath("/!root")), wants: []string{"`${GGROOT}`"}}, wantMessage: "", wantHelper: true},

		{name: "equal to first want is ok", fields: fields{localRoot: testutil.ToOSPath("/root/")}, args: args{prefix: "logprefix", got: "first", wants: []string{"first", "last"}}, wantMessage: "", wantHelper: true},
		{name: "equal to last want is ok", fields: fields{localRoot: testutil.ToOSPath("/root/")}, args: args{prefix: "logprefix", got: "last", wants: []string{"first", "last"}}, wantMessage: "", wantHelper: true},
		{name: "not equal to any wants shows last error", fields: fields{localRoot: testutil.ToOSPath("/root/")}, args: args{prefix: "logprefix", got: "neither", wants: []string{"first error", "${GGROOT last}"}}, wantMessage: fmt.Sprintf("logprefix got = %q, want = %q", "neither", testutil.ToOSPath("/root/last")), wantHelper: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := &MockEnv{
				localRoot: tt.fields.localRoot,
			}

			var r recordingT
			mock.AssertOutput(&r, tt.args.prefix, tt.args.got, tt.args.wants...)

			if tt.wantMessage != r.message {
				t.Errorf("mock.AssertOutput() message = %q, want = %q", r.message, tt.wantMessage)
			}

			if tt.wantHelper != r.helper {
				t.Errorf("mock.AssertOutput() helper = %v, want = %v", r.helper, tt.wantHelper)
			}
		})
	}
}
