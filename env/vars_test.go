package env

//spellchecker:words reflect testing github ggman internal testutil
import (
	"os"
	"reflect"
	"testing"
)

//spellchecker:words GGROOT CANFILE GGNORM

func TestReadVariables(t *testing.T) {

	// set fake environment variables for test
	defer setenv(t, "PATH", "/fake/path")()
	defer setenv(t, "HOME", "/fake/home")()
	defer setenv(t, "USERPROFILE", "/fake/home")()
	defer setenv(t, "GGROOT", "/fake/ggroot")()
	defer setenv(t, "GGMAN_CANFILE", "/fake/canfile")()
	defer setenv(t, "GGNORM", "something-fake")()

	got := ReadVariables()
	want := Variables{
		HOME:    "/fake/home",
		PATH:    "/fake/path",
		GGROOT:  "/fake/ggroot",
		CANFILE: "/fake/canfile",
		GGNORM:  "something-fake",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadVariables() = %v, want %v", got, want)
	}
}

// sets an environment variable, and returns a function to clean up
func setenv(t *testing.T, name, value string) func() {
	old := os.Getenv(name)

	// set new value
	if err := os.Setenv(name, value); err != nil {
		t.Errorf("failed to set environment variable %q: %s", name, err)
	}

	// return a function to revert to the old value
	return func() {
		if err := os.Setenv(name, old); err != nil {
			t.Errorf("failed to set environment variable %q: %s", name, err)
		}
	}
}
