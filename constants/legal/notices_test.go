//spellchecker:words legal
package legal_test

//spellchecker:words strings testing github ggman constants legal
import (
	"strings"
	"testing"

	"go.tkw01536.de/ggman/constants/legal"
)

//spellchecker:words ggman

func TestLicenses(t *testing.T) {
	t.Parallel()

	if legal.Notices == "" {
		t.Errorf("Notices is empty")
	}
	if strings.Contains(legal.Notices, "go.tkw01536.de/ggman") {
		t.Errorf("Notices contains legal information about ggman")
	}
}
