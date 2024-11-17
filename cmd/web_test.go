package cmd

//spellchecker:words path filepath testing github ggman internal mockenv
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

//spellchecker:words workdir reclone godoc

func TestCommandURL(t *testing.T) {
	mock := mockenv.NewMockEnv(t)

	clonePath := mock.Clone("git@github.com/hello/world.git", "hello", "world")

	subClonePath := filepath.Join(clonePath, "sub")
	if err := os.MkdirAll(subClonePath, os.ModePerm); err != nil {
		panic(err)
	}

	nonRepoPath := filepath.Join(clonePath, "..", "..", "example.com", "other")
	if err := os.MkdirAll(nonRepoPath, os.ModePerm); err != nil {
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
			"Open url at root",
			clonePath,
			[]string{"url"},
			0,
			"https://github.com/hello/world\n",
			"",
		},

		{
			"Print clone url at root",
			clonePath,
			[]string{"url", "--clone"},
			0,
			"git clone https://github.com/hello/world.git\n",
			"",
		},

		{
			"Print reclone url at root",
			clonePath,
			[]string{"url", "--reclone"},
			0,
			"git clone git@github.com/hello/world.git\n",
			"",
		},

		{
			"Open url at faked root",
			nonRepoPath,
			[]string{"url", "--force-repo-here"},
			0,
			"https://example.com/other\n",
			"",
		},

		{
			"Open url with custom base",
			clonePath,
			[]string{"url", "https://mybase"},
			0,
			"https://mybase/hello/world\n",
			"",
		},

		{
			"Open url with custom and prefix base",
			clonePath,
			[]string{"url", "https://mybase/", "--prefix"},
			0,
			"https://mybase/github.com/hello/world\n",
			"",
		},

		{
			"Open url with predefined base with prefix",
			clonePath,
			[]string{"url", "godoc"},
			0,
			"https://pkg.go.dev/github.com/hello/world\n",
			"",
		},

		{
			"Open url with predefined base without prefix",
			clonePath,
			[]string{"url", "travis"},
			0,
			"https://travis-ci.com/hello/world\n",
			"",
		},

		{
			"Open url with tree at root",
			clonePath,
			[]string{"url", "--tree"},
			0,
			"https://github.com/hello/world/tree/master/\n",
			"",
		},

		{
			"Do not print clone url with tree at root",
			clonePath,
			[]string{"url", "--clone", "--tree"},
			4,
			"",
			"`ggman url` does not support clone and tree arguments at the same time\n",
		},

		{
			"Do not print reclone url with tree at root",
			clonePath,
			[]string{"url", "--reclone", "--tree"},
			4,
			"",
			"`ggman url` does not support reclone and tree arguments at the same time\n",
		},

		{
			"Open url at faked root with tree",
			nonRepoPath,
			[]string{"url", "--force-repo-here", "--tree"},
			0,
			"https://example.com/other\n",
			"",
		},

		{
			"Open url with branch at root",
			clonePath,
			[]string{"url", "--branch"},
			0,
			"https://github.com/hello/world/tree/master\n",
			"",
		},

		{
			"Print clone url with branch at root",
			clonePath,
			[]string{"url", "--clone", "--branch"},
			0,
			"git clone https://github.com/hello/world.git --branch master\n",
			"",
		},

		{
			"Print reclone url with branch at root",
			clonePath,
			[]string{"url", "--reclone", "--branch"},
			0,
			"git clone git@github.com/hello/world.git --branch master\n",
			"",
		},

		{
			"Open url at faked root with branch",
			nonRepoPath,
			[]string{"url", "--force-repo-here", "--branch"},
			0,
			"https://example.com/other\n",
			"",
		},

		{
			"Open url at subpath",
			subClonePath,
			[]string{"url"},
			0,
			"https://github.com/hello/world\n",
			"",
		},

		{
			"Print clone url at subpath",
			subClonePath,
			[]string{"url", "--clone"},
			0,
			"git clone https://github.com/hello/world.git\n",
			"",
		},

		{
			"Print reclone url at subpath",
			subClonePath,
			[]string{"url", "--reclone"},
			0,
			"git clone git@github.com/hello/world.git\n",
			"",
		},

		{
			"Open url with tree at subpath",
			subClonePath,
			[]string{"url", "--tree"},
			0,
			"https://github.com/hello/world/tree/master/sub\n",
			"",
		},

		{
			"Open url with branch at subpath",
			subClonePath,
			[]string{"url", "--branch"},
			0,
			"https://github.com/hello/world/tree/master\n",
			"",
		},

		{
			"Print clone url with branch at subpath",
			subClonePath,
			[]string{"url", "--clone", "--branch"},
			0,
			"git clone https://github.com/hello/world.git --branch master\n",
			"",
		},

		{
			"Print reclone url with branch at subpath",
			subClonePath,
			[]string{"url", "--reclone", "--branch"},
			0,
			"git clone git@github.com/hello/world.git --branch master\n",
			"",
		},

		{
			"List all bases",
			clonePath,
			[]string{"url", "--list-bases"},
			0,
			"circle: https://app.circleci.com/pipelines/github\ngodoc: https://pkg.go.dev/\nlocalgodoc: http://localhost:6060/pkg/\ntravis: https://travis-ci.com\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(URL, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
