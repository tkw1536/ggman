package env

//spellchecker:words reflect testing github ggman internal testutil
import (
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/internal/testutil"
)

//spellchecker:words GGROOT CANFILE GGNORM

func TestReadVariables(t *testing.T) {
	defer testutil.MockVariables(map[string]string{
		"PATH":          "/fake/path",
		"HOME":          "/fake/home",
		"USERPROFILE":   "/fake/home",
		"GGROOT":        "/fake/ggroot",
		"GGMAN_CANFILE": "/fake/canfile",
		"GGNORM":        "something-fake",
	})()

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
