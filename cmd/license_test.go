package cmd

import (
	"fmt"
	"testing"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/constants/legal"
	"github.com/tkw1536/ggman/internal/mockenv"
)

func TestCommandLicense(t *testing.T) {
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
			"print license information",
			"",
			[]string{"license"},

			0,
			fmt.Sprintf(stringLicenseInfo, ggman.License, legal.Notices),
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(License, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
