package git

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tkw1536/ggman"
)

type gitgit struct {
	gogit
	gitPath string
}

func (gg *gitgit) Init() (err error) {
	gg.gitPath, err = gg.findgit()
	return
}

func (gg gitgit) findgit() (git string, err error) {
	// this code has been adapted from exec.LookPath in the standard library
	// it allows using a more generic path variables
	for _, git := range filepath.SplitList(gg.gitPath) {
		if git == "" { // unix shell behavior
			git = "."
		}
		git = filepath.Join(git, "git")
		d, err := os.Stat(git)
		if err != nil {
			continue
		}
		if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
			return git, nil
		}
	}
	return "", exec.ErrNotFound
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
