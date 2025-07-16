//spellchecker:words constants
package constants_test

//spellchecker:words reflect testing time ggman constants
import (
	"reflect"
	"testing"
	"time"

	"go.tkw01536.de/ggman/constants"
)

func TestBuildVersion(t *testing.T) {
	t.Parallel()

	wantBuildVersion := "v0.0.0-unknown"
	if constants.BuildVersion != wantBuildVersion {
		t.Errorf("BuildVersion = %q, want = %q", constants.BuildVersion, wantBuildVersion)
	}
}

func TestBuildTime(t *testing.T) {
	t.Parallel()

	wantBuildTime := time.Unix(0, 0).UTC()
	if !reflect.DeepEqual(constants.BuildTime, wantBuildTime) {
		t.Errorf("BuildTime = %v, want = %v", constants.BuildTime, wantBuildTime)
	}
}
