package cmd

import (
	"testing"

	"github.com/tkw1536/ggman/testutil/mockenv"
)

func TestCommandRoot(t *testing.T) {
	mock, cleanup := mockenv.NewMockEnv()
	defer cleanup()

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"show root directory",
			"",
			[]string{"root"},

			0,
			"${GGROOT}\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Root, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, stdout, tt.wantStdout, "Stdout")
			mock.AssertOutput(t, stderr, tt.wantStderr, "Stderr")
		})
	}
}
