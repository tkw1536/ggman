package constants

import (
	"reflect"
	"testing"
	"time"
)

func TestBuildVersion(t *testing.T) {
	wantBuildVersion := "v0.0.0-unknown"
	if buildVersion != wantBuildVersion {
		t.Errorf("buildVersion = %q, want = %q", buildVersion, wantBuildVersion)
	}
}

func TestBuildTime(t *testing.T) {
	wantBuildTime := time.Unix(0, 0).UTC()
	if !reflect.DeepEqual(BuildTime, wantBuildTime) {
		t.Errorf("BuildTime = %v, want = %v", BuildTime, wantBuildTime)
	}
}
