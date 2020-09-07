package git

import (
	"os/exec"

	"github.com/tkw1536/ggman"
)

type gitgit struct {
	gogit
	gitPath string
}

func (gg *gitgit) Init() (err error) {
	gg.gitPath, err = exec.LookPath("git")
	return
}

func (gg gitgit) Clone(stream ggman.IOStream, remoteURI, clonePath string, extraargs ...string) error {

	gargs := append([]string{"clone", remoteURI, clonePath}, extraargs...)

	cmd := exec.Command(gg.gitPath, gargs...)
	cmd.Stdin = stream.Stdin
	cmd.Stdout = stream.Stdout
	cmd.Stderr = stream.Stderr

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
