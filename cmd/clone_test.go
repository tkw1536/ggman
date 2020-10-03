package cmd

import (
	"testing"

	"github.com/tkw1536/ggman/testutil/mockenv"
)

func TestCommandClone(t *testing.T) {
	mock, cleanup := mockenv.NewMockEnv()
	defer cleanup()

	mock.Register("https://github.com/hello/world.git", "git@github.com:hello/world.git")

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"clone repository that doesn't exist yet",
			"",
			[]string{"clone", "https://github.com/hello/world.git"},

			0,
			"Cloning \"git@github.com:hello/world.git\" into \"${GGROOT github.com hello world}\" ...\nEnumerating objects: 3, done.\nCounting objects:  33% (1/3)\rCounting objects:  66% (2/3)\rCounting objects: 100% (3/3)\rCounting objects: 100% (3/3), done.\nTotal 3 (delta 0), reused 0 (delta 0), pack-reused 0\n",
			"",
		},

		{
			"clone existing repository",
			"",
			[]string{"clone", "https://github.com/hello/world.git"},

			1,
			"Cloning \"git@github.com:hello/world.git\" into \"${GGROOT github.com hello world}\" ...\n",
			"Unable to clone repository: Another git repository already exists in target\nlocation.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Clone, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, stdout, tt.wantStdout, "Stdout")
			mock.AssertOutput(t, stderr, tt.wantStderr, "Stderr")
		})
	}
}
