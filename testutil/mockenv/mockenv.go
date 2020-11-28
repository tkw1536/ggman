// Package mockenv contains facilities for unit testing commands
package mockenv

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	gggit "github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/testutil"
)

// MockEnv represents a new environment that can be used for testing ggman commands.
//
// The mocked environment creates a temporary folder which can be used to hold repositories.
// In order to mock a certain local state, repositories can be installed into this folder using the Register() and Install() commands.
type MockEnv struct {
	localRoot  string
	remoteRoot string

	vars     env.Variables
	plumbing *MappedPlumbing

	remoteCounter int
}

// NewMockEnv creates a new MockEnv for testing ggman programs.
// It also returns a cleanup function, and should be called as follows:
//  mock, cleanup = NewMockEnv()
//  defer cleanup()
func NewMockEnv() (*MockEnv, func()) {
	root, cleanup := testutil.TempDir()

	local := filepath.Join(root, "local")
	if err := os.Mkdir(local, os.ModePerm); err != nil {
		panic(err)
	}
	remote := filepath.Join(root, "remote")
	if err := os.Mkdir(remote, os.ModePerm); err != nil {
		panic(err)
	}

	return &MockEnv{
		localRoot:  local,
		remoteRoot: remote,

		plumbing: &MappedPlumbing{
			Plumbing: gggit.NewPlumbing(),
			URLMap:   make(map[string]string),
		},

		vars: env.Variables{
			HOME:   local,
			PATH:   "",
			GGROOT: local,
		},
	}, cleanup
}

func (mock MockEnv) resolve(path ...string) string {
	return filepath.Join(append([]string{mock.localRoot}, path...)...)
}

// Install installs the provided remote into the provided path.
// Returns the path the repository has been installed into.
// Assumes that the remote has been registered.
//
// If something goes wrong, calls panic().
func (mock *MockEnv) Install(remote string, path ...string) string {
	clonePath := mock.resolve(path...)
	err := mock.plumbing.Clone(ggman.NewNilIOStream(), remote, clonePath)
	if err != nil {
		panic(err)
	}
	return clonePath
}

// Register registers a new remote repository with the provided urls.
// All remote urls point to the same path.
// The remote repository contains one commit with the provided hash.
//
// The purpose of registering a remote is that it does not place a requirement for external services to be alive during testing.
// Calls to clone or fetch the provided repository will instead of talking to the remote talk to this dummy repository instead.
func (mock *MockEnv) Register(remotes ...string) (repo *git.Repository, hash plumbing.Hash) {
	if len(remotes) == 0 {
		panic("Register: Must provide at least one remote. ")
	}
	// generate a new fake remote path
	mock.remoteCounter++
	fakeRemotePath := filepath.Join(mock.remoteRoot, "remote"+fmt.Sprint(mock.remoteCounter))

	// create a repository
	repo = testutil.NewTestRepoAt(fakeRemotePath)
	_, hash = testutil.CommitTestFiles(repo, map[string]string{"fake.txt": remotes[0]})

	mock.plumbing.URLMap[remotes[0]] = fakeRemotePath

	// Register all the repositories.
	// Here we rely on the fact that adding "/." to the end of a path does not change the actually cloned path.
	// We can thus add it to the mapped remote, and still refer to the same repository.
	suffix := ""
	for _, remote := range remotes {
		mock.plumbing.URLMap[remote] = fakeRemotePath + suffix
		suffix += string(filepath.Separator) + "."
	}

	return repo, hash
}

// Run runs the command with the provided arguments.
// It afterwards resets the concrete value stored in command to it's zero value.
//
// The arguments should include the name of the command.
// The command is provided the given string as standard input.
//
// Commands are not executed on the real system; instead they are executed within the sandboxed environment.
// In particular all interactions with remote git repositories are intercepted, see the Register() method for details.
//
// Run returns the exit code of the provided command, along with standard output and standard error.
func (mock *MockEnv) Run(command program.Command, workdir string, stdin string, argv ...string) (code uint8, stdout, stderr string) {
	// create buffers
	stdinReader := strings.NewReader(stdin)
	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}

	// create a program and run Main()
	fakeggman := &program.Program{
		IOStream: ggman.NewIOStream(stdoutBuffer, stderrBuffer, stdinReader, 0),
	}
	fakeggman.Register(program.CloneCommand(command))

	// run the code
	err := ggman.AsError(fakeggman.Main(env.EnvironmentParameters{
		Variables: mock.vars,
		Plumbing:  mock.plumbing,
		Workdir:   workdir,
	}, argv))
	return uint8(err.ExitCode), stdoutBuffer.String(), stderrBuffer.String()
}

// regular expression used for substiution
var regexGGROOT = regexp.MustCompile(`.?\$\{GGROOT( [^\}]+)?\}.?`)

// TestingT is an interface around TestingT
type TestingT interface {
	Errorf(format string, args ...interface{})
}

// AssertOutput asserts that the standard error or output returned by Run() is equal to one of wants.
// If this is not the case, calls TestingT.Errorf() with an error message relating to the last want.
//
// For consistency across runs, strings of the form `${GGROOT a b c}` in want are resolved into an absolute path.
// Furthermore when the character preceeding the $ is a '"', additionally escapes the string.
//
// Context should be aditional information to be prefixed for the error message.
func (mock *MockEnv) AssertOutput(t TestingT, prefix, got string, wants ...string) {
	var ok bool
	var lastWant string
	for _, want := range wants {
		ok, lastWant = mock.isOutputSingle(got, want)
		if ok {
			return
		}
	}
	t.Errorf("%s got = %q, want = %q", prefix, got, lastWant)
}

// isOutputSingle normalizes want and checks if got = want.
func (mock *MockEnv) isOutputSingle(got, want string) (ok bool, normalizedWant string) {
	want = regexGGROOT.ReplaceAllStringFunc(want, func(s string) string {
		// extract the first character, actual characters, and the last character
		first := string(s[0])
		actual := s[1 : len(s)-1]
		last := string(s[len(s)-1])

		if actual[0] != '$' { // the first character was empty
			first = ""
			actual = "$" + actual
		}
		if actual[len(actual)-1] != '}' { // the last character was empty
			last = ""
			actual = actual + "}"
		}

		parts := strings.Fields(actual[:len(actual)-1])[1:] // remove trailing '}' and first part (${GGROOT)
		actual = mock.resolve(parts...)

		if first == "\"" && last == "\"" {
			return fmt.Sprintf("%q", actual)
		}
		return first + actual + last
	})
	return got == want, want
}
