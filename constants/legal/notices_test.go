//spellchecker:words legal
package legal_test

//spellchecker:words strings testing github ggman constants legal
import (
	"strings"
	"testing"

	"github.com/tkw1536/ggman/constants/legal"
)

//spellchecker:words ggman

func TestLicenses(t *testing.T) {
	t.Parallel()

	if legal.Notices == "" {
		t.Errorf("Notices is empty")
	}
	if strings.Contains(legal.Notices, "github.com/tkw1536/ggman") {
		t.Errorf("Notices contains legal information about ggman")
	}
}
