package cmd

//spellchecker:words exec runtime testing github ggman internal mockenv
import (
	"os/exec"
	"runtime"
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
)

//spellchecker:words workdir GGROOT

func setupExecTest(t *testing.T) (mock *mockenv.MockEnv) {
	mock = mockenv.NewMockEnv(t)

	mock.Clone("https://github.com/hello/world.git", "github.com", "hello", "world")
	mock.Clone("user@server.com/repo", "server.com", "user", "repo")
	mock.Clone("https://gitlab.com/hello/world.git", "gitlab.com", "hello", "world")

	return
}

func TestCommandExec_real(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows because pwd changes")
	}
	if _, err := exec.LookPath("pwd"); err != nil {
		t.Skip("pwd not found in path")
	}

	mock := setupExecTest(t)

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{

		{
			"normal exec",
			"",
			[]string{"exec", "pwd"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",
		},
		{
			"don't print repository",
			"",
			[]string{"exec", "--no-repo", "pwd"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",
			"",
		},

		{
			"be quiet",
			"",
			[]string{"exec", "--quiet", "pwd"},

			0,
			"",
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Exec, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}

func TestCommandExec_false(t *testing.T) {

	if _, err := exec.LookPath("false"); err != nil {
		t.Skip("false not found in path")
	}

	mock := setupExecTest(t)

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{

		{
			"false without force",
			"",
			[]string{"exec", "false"},

			1,
			"",
			"${GGROOT github.com hello world}\nExit status 1\n",
		},

		{
			"false with force",
			"",
			[]string{"exec", "--force", "false"},

			1,
			"",
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\nExit status 1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Exec, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}

func TestCommandExec_simulate(t *testing.T) {
	if _, err := exec.LookPath("pwd"); err != nil {
		t.Skip("pwd not found in path")
	}

	mock := setupExecTest(t)

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{

		{
			"simulate exec",
			"",
			[]string{"exec", "--simulate", "pwd"},

			0,
			"#!/bin/bash\nset -e\n\ncd `${GGROOT github.com hello world}`\necho `${GGROOT github.com hello world}`\npwd\n\ncd `${GGROOT gitlab.com hello world}`\necho `${GGROOT gitlab.com hello world}`\npwd\n\ncd `${GGROOT server.com user repo}`\necho `${GGROOT server.com user repo}`\npwd\n\n",
			"",
		},

		{
			"simulate exec with --no-repo",
			"",
			[]string{"exec", "--simulate", "--no-repo", "pwd"},

			0,
			"#!/bin/bash\nset -e\n\ncd `${GGROOT github.com hello world}`\npwd\n\ncd `${GGROOT gitlab.com hello world}`\npwd\n\ncd `${GGROOT server.com user repo}`\npwd\n\n",
			"",
		},

		{
			"simulate exec with --no-repo --force",
			"",
			[]string{"exec", "--simulate", "--no-repo", "--force", "pwd"},

			0,
			"#!/bin/bash\n\ncd `${GGROOT github.com hello world}`\npwd\n\ncd `${GGROOT gitlab.com hello world}`\npwd\n\ncd `${GGROOT server.com user repo}`\npwd\n\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Exec, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
