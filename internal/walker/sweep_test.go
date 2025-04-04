//spellchecker:words walker
package walker

//spellchecker:words path filepath reflect testing github ggman internal testutil pkglib testlib
import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/internal/testutil"
	"github.com/tkw1536/pkglib/testlib"
)

func TestSweep(t *testing.T) {
	base := testlib.TempDirAbs(t)

	// setup a directory structure for testing.
	// Make mkdir and symlink utility methods for this.

	mkdir := func(s string, files ...string) {
		t.Helper()

		path := filepath.Join(base, s)
		err := os.MkdirAll(path, os.ModePerm) // #nosec: G301 fine in a test
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			//#nosec G306
			if err := os.WriteFile(filepath.Join(path, f), nil, os.ModePerm); err != nil {
				panic(err)
			}
		}
	}

	// create a directory structure
	//
	// folders starting with f are full
	// folders starting with e are empty
	mkdir(filepath.Join("f", "f"), "file")
	mkdir(filepath.Join("f", "e"))
	mkdir(filepath.Join("e", "e1"))
	mkdir(filepath.Join("e", "e2"))

	// trimPath makes path relative to the root
	trimPath := func(path string) string {
		t, err := filepath.Rel(base, path)
		if err != nil {
			return path
		}
		return t
	}
	// trimAll makes all paths relative to the root
	trimAll := func(paths []string) {
		for idx := range paths {
			paths[idx] = trimPath(paths[idx])
		}
	}

	tests := []struct {
		name    string
		visit   SweepProcess
		params  Params
		want    []string
		wantErr bool
	}{
		{
			"sweep / without symlinks",
			nil,
			Params{
				Root: NewRealFS(base, false),
			},
			[]string{
				"e/e1",
				"e/e2",
				"f/e",
				"e",
			},
			false,
		},
		{
			"sweep /e",
			nil,
			Params{
				Root: NewRealFS(filepath.Join(base, "e"), false),
			},
			[]string{
				"e/e1",
				"e/e2",
				"e",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sweep(tt.visit, tt.params)
			trimAll(got)
			testutil.ToOSPaths(tt.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sweep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sweep() = %v, want %v", got, tt.want)
			}
		})
	}
}
