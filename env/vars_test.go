package env

import (
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/testutil"
)

func TestReadVariables(t *testing.T) {
	defer testutil.MockVariables(map[string]string{
		"PATH":          "/fake/path",
		"HOME":          "/fake/home",
		"GGROOT":        "/fake/ggroot",
		"GGMAN_CANFILE": "/fake/canfile",
	})()

	got := ReadVariables()
	want := Variables{
		HOME:    "/fake/home",
		PATH:    "/fake/path",
		GGROOT:  "/fake/ggroot",
		CANFILE: "/fake/canfile",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadVariables() = %v, want %v", got, want)
	}
}
