package cmd

import (
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

func TestCommandClone(t *testing.T) {
	mock := mockenv.NewMockEnv(t)

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
			"Cloning \"git@github.com:hello/world.git\" into \"${GGROOT github.com hello world}\" ...\n",
			"",
		},

		{
			"clone repository into local path",
			mock.Resolve(),
			[]string{"clone", "--local", "https://github.com/hello/world.git"},

			0,
			"Cloning \"git@github.com:hello/world.git\" into \"${GGROOT world}\" ...\n",
			"",
		},
		{
			"clone repository into specific path",
			mock.Resolve(),
			[]string{"clone", "--to", "somewhere", "https://github.com/hello/world.git"},

			0,
			"Cloning \"git@github.com:hello/world.git\" into \"${GGROOT somewhere}\" ...\n",
			"",
		},
		{
			"clone repository into invalid path path",
			mock.Resolve(),
			[]string{"clone", "--local", "--to", "somewhere", "https://github.com/hello/world.git"},

			4,
			"",
			"invalid destination: '--to' and '--local' may not be used together\n",
		},
		{
			"clone existing repository",
			"",
			[]string{"clone", "https://github.com/hello/world.git"},

			1,
			"Cloning \"git@github.com:hello/world.git\" into \"${GGROOT github.com hello world}\" ...\n",
			"unable to clone repository: another git repository already exists in target\nlocation\n",
		},

		{
			"clone existing repository (with force)",
			"",
			[]string{"clone", "--force", "https://github.com/hello/world.git"},

			0,
			"Cloning \"git@github.com:hello/world.git\" into \"${GGROOT github.com hello world}\" ...\nClone already exists in target location, done.\n",
			"",
		},

		{
			"clone relative path",
			"",
			[]string{"clone", "./example"},

			4,
			"",
			"invalid remote URI \"./example\": invalid scheme, not a remote path\n",
		},

		{
			"clone relative path (2)",
			"",
			[]string{"clone", "/some/example/path"},

			4,
			"",
			"invalid remote URI \"/some/example/path\": invalid scheme, not a remote path\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Clone, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
