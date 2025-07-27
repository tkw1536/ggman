//spellchecker:words constants
package ggman_test

//spellchecker:words reflect testing time ggman constants
import (
	"reflect"
	"testing"
	"time"

	"go.tkw01536.de/ggman"
)

func TestBuildVersion(t *testing.T) {
	t.Parallel()

	wantBuildVersion := "v0.0.0-unknown"
	if ggman.BuildVersion != wantBuildVersion {
		t.Errorf("BuildVersion = %q, want = %q", ggman.BuildVersion, wantBuildVersion)
	}
}

func TestBuildTime(t *testing.T) {
	t.Parallel()

	wantBuildTime := time.Unix(0, 0).UTC()
	if !reflect.DeepEqual(ggman.BuildTime, wantBuildTime) {
		t.Errorf("BuildTime = %v, want = %v", ggman.BuildTime, wantBuildTime)
	}
}
