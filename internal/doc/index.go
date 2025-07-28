package doc

//spellchecker:words strings github cobra
import (
	"strings"

	"github.com/spf13/cobra"
)

// IndexFilenames returns a map from filename to command path.
func IndexFilenames(cmd *cobra.Command, extension string) map[string]string {
	mp := make(map[string]string, CountCommands(cmd))

	for cmd := range AllCommands(cmd) {
		path := cmd.CommandPath()

		name := strings.ReplaceAll(path, " ", "_") + "." + extension
		mp[name] = path
	}

	return mp
}
