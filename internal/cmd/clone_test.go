package cmd_test

//spellchecker:words testing ggman internal mockenv
import (
	"testing"

	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words GGROOT tparallel paralleltest

//nolint:tparallel,paralleltest
func TestCommandClone(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	mock.Register("https://github.com/hello/world.git", "git@github.com:hello/world.git")
	mock.Register("https://github.com/hello/world2.git", "git@github.com:hello/world2.git")
	mock.Register("https://github.com/hello/world3.git")
	mock.Register("https://github.com/hello/world4.git", "git@github.com:hello/world4.git")
	mock.Register("https://github.com/hello/world5.git", "git@github.com:hello/world5.git")

	// These tests should not be run in parallel, but treated as a single linear test.
	// Each test case depends on the previous one and implicitly relies on the fact that
	// the previous test left the environment in a specific state.
	tests := []struct {
		name    string
		workDir string
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
			[]string{"clone", "--plain", "https://github.com/hello/world.git"},

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
			[]string{"clone", "--plain", "--to", "somewhere", "https://github.com/hello/world.git"},

			4,
			"",
			"invalid destination: \"--to\" and \"--plain\" may not be used together\n",
		},
		{
			"clone existing repository",
			"",
			[]string{"clone", "https://github.com/hello/world.git"},

			1,
			"Cloning \"git@github.com:hello/world.git\" into \"${GGROOT github.com hello world}\" ...\n",
			"unable to clone repository: another git repository already exists in target location\n",
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
			"\"./example\": invalid remote URI: invalid scheme, not a remote path\n",
		},

		{
			"clone relative path (2)",
			"",
			[]string{"clone", "--", "/some/example/path"},

			4,
			"",
			"\"/some/example/path\": invalid remote URI: invalid scheme, not a remote path\n",
		},

		{
			"clone repository with exact url",
			"",
			[]string{"clone", "--exact-url", "https://github.com/hello/world3.git"},

			0,
			"Cloning \"https://github.com/hello/world3.git\" into \"${GGROOT github.com hello world3}\" ...\n",
			"",
		},

		{
			"clone repository and args",
			// this doesn't actually clone (because we don't have a real git)
			// but at least parses the args
			"",
			[]string{"clone", "https://github.com/hello/world4.git", "--", "--depth", "1"},

			1,
			"Cloning \"git@github.com:hello/world4.git\" into \"${GGROOT github.com hello world4}\" ...\n",
			"external `git` not found, can not pass any additional arguments to `git clone`: --depth 1\n",
		},

		{
			"clone existing repository (with overwrite)",
			"",
			[]string{"clone", "--overwrite", "https://github.com/hello/world.git"},

			0,
			"Deleting existing directory \"${GGROOT github.com hello world}\"\nCloning \"git@github.com:hello/world.git\" into \"${GGROOT github.com hello world}\" ...\n",
			"",
		},

		{
			"clone non-existing repository (with overwrite)",
			"",
			[]string{"clone", "--overwrite", "https://github.com/hello/world5.git"},

			0,
			"Cloning \"git@github.com:hello/world5.git\" into \"${GGROOT github.com hello world5}\" ...\n",
			"",
		},

		{
			"fail to clone existing repository (it's still there)",
			"",
			[]string{"clone", "https://github.com/hello/world.git"},

			1,
			"Cloning \"git@github.com:hello/world.git\" into \"${GGROOT github.com hello world}\" ...\n",
			"unable to clone repository: another git repository already exists in target location\n",
		},

		{
			"clone repository with overwrite and force",
			mock.Resolve(),
			[]string{"clone", "--force", "--overwrite", "https://github.com/hello/world.git"},

			4,
			"",
			"\"--overwrite\" and \"--force\" are incompatible\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, tt.workDir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
