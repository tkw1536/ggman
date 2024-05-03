//spellchecker:words legal
package legal

//spellchecker:words strings testing
import (
	"strings"
	"testing"
)

//spellchecker:words ggman

func TestLicenses(t *testing.T) {
	if Notices == "" {
		t.Errorf("Notices is empty")
	}
	if strings.Contains(Notices, "github.com/tkw1536/ggman") {
		t.Errorf("Notices contains legal information about ggman")
	}
}
