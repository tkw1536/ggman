package legal

import (
	"strings"
	"testing"
)

func TestLicenses(t *testing.T) {
	if Notices == "" {
		t.Errorf("Notices is empty")
	}
	if strings.Contains(Notices, "github.com/tkw1536/ggman") {
		t.Errorf("Notices contains legal information about ggman")
	}
}
