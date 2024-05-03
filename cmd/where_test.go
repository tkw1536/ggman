package cmd

//spellchecker:words testing github ggman internal mockenv
import (
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

//spellchecker:words ggman GGROOT workdir

func TestCommandWhere(t *testing.T) {
	mock := mockenv.NewMockEnv(t)

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"show directory of repository",
			"",
			[]string{"where", "https://github.com/hello/world.git"},

			0,
			"${GGROOT github.com hello world}\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Where, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
