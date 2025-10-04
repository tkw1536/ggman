package cmd_test

//spellchecker:words testing github config ggman internal mockenv
import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words GGROOT workdir tparallel paralleltest

//nolint:tparallel,paralleltest
func TestCommandFix(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	mock.Register("https://github.com/hello/world.git", "git@github.com:hello/world.git")
	mock.Install(t.Context(), "https://github.com/hello/world.git", "github.com", "hello", "world")

	mock.Register("user@server.com/repo", "git@server.com:user/repo.git")
	mock.Install(t.Context(), "user@server.com/repo", "server.com", "user", "repo")

	mock.Register("https://gitlab.com/hello/world.git", "git@gitlab.com:hello/world.git")
	mock.Install(t.Context(), "https://gitlab.com/hello/world.git", "gitlab.com", "hello", "world")

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"simulate fixing remotes of all repositories",
			"",
			[]string{"fix", "--simulate"},

			0,
			"Simulate fixing remote of \"${GGROOT github.com hello world}\"\nUpdating origin: https://github.com/hello/world.git -> git@github.com:hello/world.git\nSimulate fixing remote of \"${GGROOT gitlab.com hello world}\"\nUpdating origin: https://gitlab.com/hello/world.git -> git@gitlab.com:hello/world.git\nSimulate fixing remote of \"${GGROOT server.com user repo}\"\nUpdating origin: user@server.com/repo -> git@server.com:user/repo.git\n",
			"",
		},

		{
			"actually fixing remotes of all repositories",
			"",
			[]string{"fix"},

			0,
			"Fixing remote of \"${GGROOT github.com hello world}\"\nUpdating origin: https://github.com/hello/world.git -> git@github.com:hello/world.git\nFixing remote of \"${GGROOT gitlab.com hello world}\"\nUpdating origin: https://gitlab.com/hello/world.git -> git@gitlab.com:hello/world.git\nFixing remote of \"${GGROOT server.com user repo}\"\nUpdating origin: user@server.com/repo -> git@server.com:user/repo.git\n",
			"",
		},

		{
			"fixing remotes of fixed repositories",
			"",
			[]string{"fix"},

			0,
			"",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}

//nolint:tparallel,paralleltest
func TestCommandFix_Prune(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	mock.Register("git@github.com:hello/world.git")
	_, remotes := mock.Register("git@github.com:hello/world2.git")

	repoPath := mock.Install(t.Context(), "git@github.com:hello/world.git", "github.com", "hello", "world")

	// Add upstream remote to the repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		panic(err)
	}
	if _, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "upstream",
		URLs: []string{remotes[0]}, // yes, it points to itself - but that should never be touched anyways
	}); err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			name:    "nothing gets pruned when --prune-remotes is not set",
			workdir: "",
			args:    []string{"fix", "--simulate"},

			wantStdout: "",
			wantCode:   0,
		},
		{
			name:    "simulate pruning remotes of all repositories",
			workdir: "",
			args:    []string{"fix", "--simulate", "--prune-remotes"},

			wantStdout: "Found unused remote \"upstream\" in \"${GGROOT github.com hello world}\"\n",
			wantCode:   0,
		},
		{
			name:    "actually pruning remotes of all repositories",
			workdir: "",
			args:    []string{"fix", "--prune-remotes"},

			wantStdout: "Removing unused remote \"upstream\" from \"${GGROOT github.com hello world}\"\n",
			wantCode:   0,
		},
		{
			name:    "no more remotes to prune",
			workdir: "",
			args:    []string{"fix", "--prune-remotes"},

			wantStdout: "",
			wantCode:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
