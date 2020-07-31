package gitwrap

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

func (gg *gitgit) Clone(remoteURI, clonePath string, extraargs ...string) (code int, err error) {

	gargs := append([]string{"clone", remoteURI, clonePath}, extraargs...)

	cmd := exec.Command(gg.gitPath, gargs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run the command and don't treat exitError as an error
	err = cmd.Run()
	if exitError, isExitError := err.(*exec.ExitError); isExitError {
		code = exitError.ExitCode()
		err = nil
	}
	return
}

func init() {
	// check that goGitImpl is a git implementation
	var _ GitImplementation = (*gitgit)(nil)
}
