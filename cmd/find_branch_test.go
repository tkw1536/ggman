package cmd_test

//spellchecker:words testing github config ggman internal mockenv
import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/tkw1536/ggman/cmd"
	"github.com/tkw1536/ggman/internal/mockenv"
)

//spellchecker:words GGROOT workdir

func TestCommandFindBranch(t *testing.T) {
	mock := mockenv.NewMockEnv(t)

	// with branch 'branch'
	clonePath := mock.Clone("https://github.com/hello/world.git", "github.com", "hello", "world")
	repo, err := git.PlainOpen(clonePath)
	if err != nil {
		panic(err)
	}
	if err := repo.CreateBranch(&config.Branch{Name: "branch"}); err != nil {
		panic(err)
	}

	// with branch 'branch'
	clonePath = mock.Clone("user@server.com/repo", "server.com", "user", "repo")
	repo, err = git.PlainOpen(clonePath)
	if err != nil {
		panic(err)
	}
	if err := repo.CreateBranch(&config.Branch{Name: "branch"}); err != nil {
		panic(err)
	}

	// with only master branch
	repo = mock.Register("https://gitlab.com/hello/world.git")
	if err := repo.CreateBranch(&config.Branch{Name: "branchC"}); err != nil {
		panic(err)
	}
	mock.Install("https://gitlab.com/hello/world.git", "gitlab.com", "hello", "world")

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"find master branches",
			"",
			[]string{"find-branch", "master"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",
			"",
		},

		{
			"find 'branch' branches",
			"",
			[]string{"find-branch", "branch"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT server.com user repo}\n",
			"",
		},

		{
			"find 'branch' branches with --exit-code",
			"",
			[]string{"find-branch", "branch", "--exit-code"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT server.com user repo}\n",
			"",
		},

		{
			"find 'fake' branches",
			"",
			[]string{"find-branch", "fake"},

			0,
			"",
			"",
		},

		{
			"find 'fake' branches with --exit-code",
			"",
			[]string{"find-branch", "fake", "--exit-code"},

			1,
			"",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(cmd.FindBranch, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
