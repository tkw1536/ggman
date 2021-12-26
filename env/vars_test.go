package env

import (
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/internal/testutil"
)

func TestReadVariables(t *testing.T) {
	defer testutil.MockVariables(map[string]string{
		"PATH":          "/fake/path",
		"HOME":          "/fake/home",
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
