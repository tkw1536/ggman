package git

import (
	"os"
	"os/exec"
)

type gitgit struct {
	gogit
	gitPath string
}

func (gg *gitgit) Init() (err error) {
	gg.gitPath, err = exec.LookPath("git")
	return
}

func (gg *gitgit) Clone(remoteURI, clonePath string, extraargs ...string) error {

	gargs := append([]string{"clone", remoteURI, clonePath}, extraargs...)

	cmd := exec.Command(gg.gitPath, gargs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run the underlying command, but treat ExitError specially by turning it into a ExitError
	err := cmd.Run()
	if exitError, isExitError := err.(*exec.ExitError); isExitError {
		err = ExitError{error: err, Code: exitError.ExitCode()}
	}
	return err
}

func init() {
	// check that goGitImpl is a git implementation
	var _ Plumbing = (*gitgit)(nil)
}
