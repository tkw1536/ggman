// Package mockenv contains facilities for unit testing commands
package mockenv

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/alessio/shellescape"
	"github.com/go-git/go-git/v5"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	gggit "github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/internal/path"
	"github.com/tkw1536/ggman/internal/testutil"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/reflectx"
	"github.com/tkw1536/pkglib/stream"
	"github.com/tkw1536/pkglib/testlib"
)

// TODO: Consider generalizing this part of the test suite.
// so that we can re-use it for other programs too.

// MockEnv represents a new environment that can be used for testing ggman commands.
//
// The mocked environment creates a temporary folder which can be used to hold repositories.
// In order to mock a certain local state, repositories can be installed into this folder using the Register() and Install() commands.
type MockEnv struct {
	localRoot  string
	remoteRoot string

	vars     env.Variables
	plumbing DevPlumbing

	remoteCounter int
}

// NewMockEnv creates a new MockEnv for testing ggman programs.
func NewMockEnv(t *testing.T) *MockEnv {
	root := testlib.TempDirAbs(t)

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

		plumbing: DevPlumbing{
			Plumbing: gggit.NewPlumbing(),

			SilenceStderr: true,
			URLMap:        make(map[string]string),
		},

		vars: env.Variables{
			HOME:   local,
			PATH:   "",
			GGROOT: local,
		},
	}
}

// Resolve resolves a local path within this environment
func (mock MockEnv) Resolve(path ...string) string {
	return filepath.Join(append([]string{mock.localRoot}, path...)...)
}

// Install installs the provided remote into the provided path.
// Returns the path the repository has been installed into.
// Assumes that the remote has been registered.
//
// When the remote has not been registered, consider using Install instead.
//
// If something goes wrong, calls panic().
func (mock *MockEnv) Install(remote string, path ...string) string {
	clonePath := mock.Resolve(path...)
	err := mock.plumbing.Clone(stream.FromNil(), remote, clonePath)
	if err != nil {
		panic(err)
	}
	return clonePath
}

// Clone is like Install, but calls Register(remote) beforehand.
// Returns the return value of Install.
//
// This function is untested because Register and Install are tested.
func (mock *MockEnv) Clone(remote string, path ...string) (clonePath string) {
	mock.Register(remote)
	return mock.Install(remote, path...)
}

// Register registers a new remote repository with the provided urls.
// All remote urls point to the same path.
// Returns a reference to the remote repository.
//
// The purpose of registering a remote is that it does not place a requirement for external services to be alive during testing.
// Calls to clone or fetch the provided repository will instead of talking to the remote talk to this dummy repository instead.
func (mock *MockEnv) Register(remotes ...string) (repo *git.Repository) {
	if len(remotes) == 0 {
		panic("Register: Must provide at least one remote. ")
	}
	// generate a new fake remote path
	mock.remoteCounter++
	fakeRemotePath := filepath.Join(mock.remoteRoot, "remote"+fmt.Sprint(mock.remoteCounter))

	// create a repository
	repo = testutil.NewTestRepoAt(fakeRemotePath, "")
	testutil.CommitTestFiles(repo, map[string]string{"fake.txt": remotes[0]})

	// Register all the repositories.
	// Here we rely on the fact that adding "/." to the end of a path does not change the actually cloned path.
	// We can thus add it to the mapped remote, and still refer to the same repository.
	suffix := ""
	for _, remote := range remotes {
		mock.plumbing.URLMap[remote] = fakeRemotePath + suffix
		suffix += path.Separator + "."
	}

	return repo
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
func (mock *MockEnv) Run(command ggman.Command, workdir string, stdin string, argv ...string) (code uint8, stdout, stderr string) {
	// create buffers
	stdinReader := strings.NewReader(stdin)
	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}

	// create a program and run Main()
	fakeggman := ggman.NewProgram()
	fakeggman.Register(reflectx.Copy(command, true))

	stream := stream.NewIOStream(stdoutBuffer, stderrBuffer, stdinReader, 0)

	// run the code
	err := exit.AsError(fakeggman.Main(stream, env.Parameters{
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
	Errorf(format string, args ...any)
}

// AssertOutput asserts that the standard error or output returned by Run() is equal to one of wants.
// If this is not the case, calls TestingT.Errorf() with an error message relating to the last want.
//
// For consistency across runs, strings of the form `${GGROOT a b c}` in want are resolved into an absolute path.
// Furthermore when `${}` is surrounded by "s, (e.g. "${GGROOT a b c}"), go quotes the string.
// When text is instead surrounded by “s`s, (e.g. `${GGROOT a b c}`) shell escapes the string.
//
// Context should be aditional information to be prefixed for the error message.
func (mock *MockEnv) AssertOutput(t TestingT, prefix, got string, wants ...string) {
	var lastWant string
	for _, want := range wants {
		lastWant = mock.interpolate(want)
		if lastWant == got {
			return
		}
	}
	t.Errorf("%s got = %q, want = %q", prefix, got, lastWant)
}

// interpolate interpolates the striing values by replacing all ins
func (mock *MockEnv) interpolate(value string) (result string) {
	return regexGGROOT.ReplaceAllStringFunc(value, func(s string) string {
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
		actual = mock.Resolve(parts...)

		if first == "\"" && last == "\"" {
			return fmt.Sprintf("%q", actual)
		}
		if first == "`" && last == "`" {
			return shellescape.Quote(actual)
		}
		return first + actual + last
	})
}
