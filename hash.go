//spellchecker:words ggman
package ggman

//spellchecker:words crypto encoding errors pkglib errorsx
import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"go.tkw01536.de/pkglib/errorsx"
)

//spellchecker:words nosec

var (
	errNoExecutable         = errors.New("failed to find executable")
	errOpenExecutableFailed = errors.New("failed to open executable")
	errFailedToComputeHash  = errors.New("failed to compute hash")
)

// BuildHash returns the sha256 hash of the current executable.
// If the hash cannot be computed, it returns an empty string.
func BuildHash() (s string, e error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("%w: %w", errNoExecutable, err)
	}

	f, err := os.Open(exe) /* #nosec G304 -- open the current executable */
	if err != nil {
		return "", fmt.Errorf("%w: %w", errOpenExecutableFailed, err)
	}
	defer errorsx.Close(f, &e, "failed to close executable")

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("%w: %w", errFailedToComputeHash, err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
