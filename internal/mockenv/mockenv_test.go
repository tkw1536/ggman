// Package mockenv contains facilities for unit testing commands
//
//spellchecker:words mockenv
package mockenv_test

//spellchecker:words path filepath testing ggman internal mockenv pkglib testlib
import (
	"fmt"
	"path/filepath"
	"testing"

	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/ggman/internal/mockenv"
	"go.tkw01536.de/pkglib/testlib"
)

//spellchecker:words GGROOT logprefix

// mockEnvRunCommand.
type mockEnvRunCommand struct {
	Positional struct {
		Argv []string
	} `positional-args:"true"`
}

func (mockEnvRunCommand) Description() ggman.Description {
	return ggman.Description{
		Command: "fake",
		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}
func (mockEnvRunCommand) AfterParse() error { return nil }
func (me mockEnvRunCommand) Run(context ggman.Context) error {
	clonePath := filepath.Join(context.Environment.Root, "server.com", "repo")
	remote, _ := context.Environment.Git.GetRemote(clonePath, "")

	if _, err := fmt.Fprintf(context.Stdout, "path=%s remote=%s\n", clonePath, remote); err != nil {
		return fmt.Errorf("failed to write path and remote: %w", err)
	}
	if _, err := fmt.Fprintf(context.Stderr, "got args: %v\n", me.Positional.Argv); err != nil {
		return fmt.Errorf("failed write arguments: %w", err)
	}

	return nil
}

func TestMockEnv_RunLegacy(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	// create a fake repository and install it into the mock
	repo := "https://server.com:repo"
	mock.Register(repo)
	clonePath := mock.Install(repo, "server.com", "repo")

	cmd := ggman.Command(&mockEnvRunCommand{})

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
			t.Parallel()

			gotCode, gotStdout, gotStderr := mock.RunLegacy(cmd, "", "", append([]string{"fake"}, tt.args...)...)
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
