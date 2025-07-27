//spellchecker:words legal
package ggman_test

//spellchecker:words strings testing ggman constants legal
import (
	"strings"
	"testing"

	"go.tkw01536.de/ggman"
)

//spellchecker:words ggman

func TestLicenses(t *testing.T) {
	t.Parallel()

	if ggman.Notices == "" {
		t.Errorf("Notices is empty")
	}
	if strings.Contains(ggman.Notices, "go.tkw01536.de/ggman") {
		t.Errorf("Notices contains legal information about ggman")
	}
}
