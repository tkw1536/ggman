// Package mockenv contains facilities for unit testing commands
//
//spellchecker:words mockenv
package mockenv

//spellchecker:words bytes context path filepath regexp strconv strings testing essio shellescape github cobra ggman gggit internal dirs testutil pkglib exit stream testlib
import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"al.essio.dev/pkg/shellescape"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/env"
	gggit "go.tkw01536.de/ggman/git"
	"go.tkw01536.de/ggman/internal/dirs"
	"go.tkw01536.de/ggman/internal/testutil"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/stream"
	"go.tkw01536.de/pkglib/testlib"
)

//spellchecker:words GGROOT workdir sandboxed

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
	t.Helper()

	root := testlib.TempDirAbs(t)

	local := filepath.Join(root, "local")
	if err := os.Mkdir(local, dirs.NewModBits); err != nil {
		panic(err)
	}
	remote := filepath.Join(root, "remote")
	if err := os.Mkdir(remote, dirs.NewModBits); err != nil {
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

// Resolve resolves a local path within this environment.
func (mock *MockEnv) Resolve(path ...string) string {
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
// Remotes must not have been registered before, or panic() will be called.
//
// The purpose of registering a remote is that it does not place a requirement for external services to be alive during testing.
// Calls to clone or fetch the provided repository will instead of talking to the remote talk to this dummy repository instead.
func (mock *MockEnv) Register(remotes ...string) (repo *git.Repository) {
	if len(remotes) == 0 {
		panic("Register: Must provide at least one remote. ")
	}

	// check that we have fresh remotes for all urls
	for _, remote := range remotes {
		if _, ok := mock.plumbing.URLMap[remote]; ok {
			panic("Register: remote " + remote + " already registered")
		}
	}

	// generate a new fake remote path
	mock.remoteCounter++
	fakeRemotePath := filepath.Join(mock.remoteRoot, "remote"+strconv.Itoa(mock.remoteCounter))

	// create a repository
	repo = testutil.NewTestRepoAt(fakeRemotePath, "")
	testutil.CommitTestFiles(repo, map[string]string{"fake.txt": remotes[0]})

	// Register all the repositories.
	// Here we rely on the fact that adding "/." to the end of a path does not change the actually cloned path.
	// We can thus add it to the mapped remote, and still refer to the same repository.
	suffix := ""
	for _, remote := range remotes {
		mock.plumbing.URLMap[remote] = fakeRemotePath + suffix
		suffix += string(os.PathSeparator) + "."
	}

	return repo
}

// Run runs the command with the provided arguments.
// cmdFactory should be a factory function to create the root command, typically [cmd.NewCommand].
// It afterwards resets the concrete value stored in command to it's zero value.
//
// The arguments should include the name of the command.
// The command is provided the given string as standard input.
//
// Commands are not executed on the real system; instead they are executed within the sandboxed environment.
// In particular all interactions with remote git repositories are intercepted, see the Register() method for details.
//
// It returns the exit code of the provided command, along with standard output and standard error.
func (mock *MockEnv) Run(t *testing.T, cmdFactory func(context.Context, env.Parameters) *cobra.Command, workdir string, stdin string, argv ...string) (code uint8, stdout, stderr string) {
	t.Helper()

	stdinReader := strings.NewReader(stdin)
	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}

	fake := cmdFactory(t.Context(), env.Parameters{
		Variables: mock.vars,
		Plumbing:  mock.plumbing,
		Workdir:   workdir,
	})
	fake.SetIn(stdinReader)
	fake.SetOut(stdoutBuffer)
	fake.SetErr(stderrBuffer)

	fake.SetArgs(argv)
	cmdErr := fake.Execute()
	exitCode, _ := exit.CodeFromError(cmdErr)
	if exitCode != 0 {
		errStr := fmt.Sprint(cmdErr)
		if errStr != "" {
			fmt.Fprintln(stderrBuffer, cmdErr)
		}
	}

	return uint8(exitCode), stdoutBuffer.String(), stderrBuffer.String()
}

// regular expression used for substitution.
var regexGGROOT = regexp.MustCompile(`.?\$\{GGROOT( [^\}]+)?\}.?`)

// TestingT is an interface around TestingT.
type TestingT interface {
	Errorf(format string, args ...any)
	Helper()
}

// AssertOutput asserts that the standard error or output returned by Run() is equal to one of wants.
// If this is not the case, calls TestingT.Errorf() with an error message relating to the last want.
//
// For consistency across runs, strings of the form `${GGROOT a b c}` in want are resolved into an absolute path.
// Furthermore when `${}` is surrounded by "s, (e.g. "${GGROOT a b c}"), go quotes the string.
// When text is instead surrounded by “s`s, (e.g. `${GGROOT a b c}`) shell escapes the string.
//
// Context should be additional information to be prefixed for the error message.
func (mock *MockEnv) AssertOutput(t testutil.TestingT, prefix, got string, wants ...string) {
	t.Helper()

	var lastWant string
	for _, want := range wants {
		lastWant = mock.interpolate(want)
		if lastWant == got {
			return
		}
	}
	t.Errorf("%s got = %q, want = %q", prefix, got, lastWant)
}

// interpolate interpolates the string values by replacing all ins.
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
			actual += "}"
		}

		parts := strings.Fields(actual[:len(actual)-1])[1:] // remove trailing '}' and first part (${GGROOT)
		actual = mock.Resolve(parts...)

		if first == "\"" && last == "\"" {
			return strconv.Quote(actual)
		}
		if first == "`" && last == "`" {
			return shellescape.Quote(actual)
		}
		return first + actual + last
	})
}
