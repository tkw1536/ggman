package cmd_test

//spellchecker:words testing ggman constants legal internal mockenv
import (
	"fmt"
	"testing"

	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/cmd"
	"go.tkw01536.de/ggman/constants/legal"
	"go.tkw01536.de/ggman/internal/mockenv"
)

//spellchecker:words testing ggman constants legal internal mockenv

func TestCommandLicense(t *testing.T) {
	t.Parallel()

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
			"print license information",
			"",
			[]string{"license"},

			0,
			fmt.Sprintf(cmd.StringLicenseInfo, ggman.License, legal.Notices),
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, stdout, stderr := mock.Run(t, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
