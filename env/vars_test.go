package env_test

//spellchecker:words reflect testing github ggman
import (
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/env"
)

//spellchecker:words GGROOT CANFILE GGNORM USERPROFILE GGMAN

func TestReadVariables(t *testing.T) {
	// set fake environment variables for test
	t.Setenv("PATH", "/fake/path")
	t.Setenv("HOME", "/fake/home")
	t.Setenv("USERPROFILE", "/fake/home")
	t.Setenv("GGROOT", "/fake/ggroot")
	t.Setenv("GGMAN_CANFILE", "/fake/canfile")
	t.Setenv("GGNORM", "something-fake")

	got := env.ReadVariables()
	want := env.Variables{
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
