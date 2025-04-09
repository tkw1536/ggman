package env

//spellchecker:words reflect testing
import (
	"reflect"
	"testing"
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
