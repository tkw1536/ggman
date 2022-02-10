// Package mockenv contains facilities for unit testing commands
package mockenv

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/internal/path"
)

// recordingT records a message passed to Errorf()
type recordingT struct {
	message string
}

func (f *recordingT) Errorf(format string, args ...interface{}) {
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
		{"no replacement equal", fields{path.ToOSPath("/root/")}, args{"logprefix", "example", []string{"example"}}, ""},
		{"no replacement not equal", fields{path.ToOSPath("/root/")}, args{"logprefix", "example", []string{"example2"}}, "logprefix got = \"example\", want = \"example2\""},

		{"replace only ggroot ok", fields{path.ToOSPath("/root/")}, args{"logprefix", "prefix " + path.ToOSPath("/root") + " suffix", []string{"prefix ${GGROOT} suffix"}}, ""},
		{"replace only ggroot not ok", fields{path.ToOSPath("/root/")}, args{"logprefix", "prefix " + path.ToOSPath("/root") + " suffix", []string{"prefix ${GGROOT}/sub suffix"}}, fmt.Sprintf("logprefix got = %q, want = %q", "prefix "+path.ToOSPath("/root")+" suffix", "prefix "+path.ToOSPath("/root")+"/sub suffix")},

		{"replace full path ok", fields{path.ToOSPath("/root/")}, args{"logprefix", "prefix " + path.ToOSPath("/root/a/b") + " suffix", []string{"prefix ${GGROOT a b} suffix"}}, ""},
		{"replace full path not ok", fields{path.ToOSPath("/root/")}, args{"logprefix", "prefix " + path.ToOSPath("/root") + " suffix", []string{"prefix ${GGROOT a b} suffix"}}, fmt.Sprintf("logprefix got = %q, want = %q", "prefix "+path.ToOSPath("/root")+" suffix", "prefix "+path.ToOSPath("/root/a/b")+" suffix")},

		{"escape path with quotes", fields{path.ToOSPath("/root/")}, args{"logprefix", fmt.Sprintf("%q", path.ToOSPath("/root")), []string{"\"${GGROOT}\""}}, ""},
		{"escape path with `s", fields{path.ToOSPath("/!root/")}, args{"logprefix", shellescape.Quote(path.ToOSPath("/!root")), []string{"`${GGROOT}`"}}, ""},

		{"equal to first want is ok", fields{path.ToOSPath("/root/")}, args{"logprefix", "first", []string{"first", "last"}}, ""},
		{"equal to last want is ok", fields{path.ToOSPath("/root/")}, args{"logprefix", "last", []string{"first", "last"}}, ""},
		{"not equal to any wants shows last error", fields{path.ToOSPath("/root/")}, args{"logprefix", "neither", []string{"first error", "${GGROOT last}"}}, fmt.Sprintf("logprefix got = %q, want = %q", "neither", path.ToOSPath("/root/last"))},
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

// mockEnvRunCommand
type mockEnvRunCommand struct{}

func (mockEnvRunCommand) BeforeRegister(program *ggman.Program) {}

func (mockEnvRunCommand) Description() ggman.Description {
	return ggman.Description{
		Name:       "fake",
		PosArgsMax: -1,
		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}
func (mockEnvRunCommand) AfterParse() error { return nil }
func (mockEnvRunCommand) Run(context ggman.Context) error {
	clonePath := filepath.Join(context.Runtime().Root, "server.com", "repo")
	remote, _ := context.Runtime().Git.GetRemote(clonePath)

	fmt.Fprintf(context.Stdout, "path=%s remote=%s\n", clonePath, remote)
	fmt.Fprintf(context.Stderr, "got args: %v\n", context.Args.Arguments.Pos)

	return nil
}

func TestMockEnv_Run(t *testing.T) {
	mock := NewMockEnv(t)

	// create a fake repository and install it into the mock
	repo := "https://server.com:repo"
	mock.Register(repo)
	clonePath := mock.Install(repo, "server.com", "repo")

	cmd := ggman.Command(mockEnvRunCommand{})

	tests := []struct {
		name       string
		args       []string
		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"simple args",
			[]string{"a", "b", "c"},
			0,
			fmt.Sprintf("path=%s remote=%s\n", clonePath, repo),
			"got args: [a b c]\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotStdout, gotStderr := mock.Run(cmd, "", "", append([]string{"fake"}, tt.args...)...)
			if gotCode != tt.wantCode {
				t.Errorf("MockEnv.Run() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
			if gotStdout != tt.wantStdout {
				t.Errorf("MockEnv.Run() gotStdout = %v, want %v", gotStdout, tt.wantStdout)
			}
			if gotStderr != tt.wantStderr {
				t.Errorf("MockEnv.Run() gotStderr = %v, want %v", gotStderr, tt.wantStderr)
			}
		})
	}

}
