//spellchecker:words constants
package constants_test

//spellchecker:words reflect testing time github ggman constants
import (
	"reflect"
	"testing"
	"time"

	"github.com/tkw1536/ggman/constants"
)

func TestBuildVersion(t *testing.T) {
	wantBuildVersion := "v0.0.0-unknown"
	if constants.BuildVersion != wantBuildVersion {
		t.Errorf("BuildVersion = %q, want = %q", constants.BuildVersion, wantBuildVersion)
	}
}

func TestBuildTime(t *testing.T) {
	wantBuildTime := time.Unix(0, 0).UTC()
	if !reflect.DeepEqual(constants.BuildTime, wantBuildTime) {
		t.Errorf("BuildTime = %v, want = %v", constants.BuildTime, wantBuildTime)
	}
}
