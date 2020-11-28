package cmd

import (
	"testing"

	"github.com/tkw1536/ggman/testutil/mockenv"
)

func TestCommandRelocate(t *testing.T) {
	mock, cleanup := mockenv.NewMockEnv()
	defer cleanup()

	mock.Register("https://github.com/right/directory.git")
	mock.Install("https://github.com/right/directory.git", "github.com", "right", "directory")

	mock.Register("https://github.com/correct/directory.git")
	mock.Install("https://github.com/correct/directory.git", "github.com", "incorrect", "directory")

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"relocate with simulate",
			"",
			[]string{"relocate", "--simulate"},

			0,
			"mkdir -p ${GGROOT github.com correct}\nmv ${GGROOT github.com incorrect directory} ${GGROOT github.com correct directory}\n",

			"",
		},

		{
			"relocate without simulate",
			"",
			[]string{"relocate"},

			0,
			"mkdir -p ${GGROOT github.com correct}\nmv ${GGROOT github.com incorrect directory} ${GGROOT github.com correct directory}\n",

			"",
		},

		{
			"nothing to relocate",
			"",
			[]string{"relocate"},

			0,
			"",

			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Relocate, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
