//spellchecker:words walker
package walker_test

//spellchecker:words path filepath reflect testing ggman internal testutil walker pkglib testlib
import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"go.tkw01536.de/ggman/internal/testutil"
	"go.tkw01536.de/ggman/internal/walker"
	"go.tkw01536.de/pkglib/testlib"
)

func TestScan(t *testing.T) {
	t.Parallel()

	base := testlib.TempDirAbs(t)

	// setup a directory structure for testing.
	// Make mkdir and symlink utility methods for this.

	mkdir := func(s string) {
		err := os.MkdirAll(filepath.Join(base, s), 0750)
		if err != nil {
			panic(err)
		}
	}

	symlink := func(oldName, newName string) {
		err := os.Symlink(filepath.Join(base, oldName), filepath.Join(base, newName))
		if err != nil {
			panic(err)
		}
	}

	mkdir(filepath.Join("a", "aa", "aaa"))
	mkdir(filepath.Join("a", "aa", "aab"))
	mkdir(filepath.Join("a", "aa", "aac"))
	mkdir(filepath.Join("a", "ab", "aba"))
	mkdir(filepath.Join("a", "ab", "abb"))
	mkdir(filepath.Join("a", "ab", "abc"))
	mkdir(filepath.Join("a", "ac", "aca"))
	mkdir(filepath.Join("a", "ac", "acb"))
	mkdir(filepath.Join("a", "ac", "acc"))

	symlink("", filepath.Join("a", "aa", "linked"))

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
		visit   walker.ScanProcess
		params  walker.Params
		want    []string
		wantErr bool
	}{
		{
			"scan /",
			nil,
			walker.Params{
				Root: walker.NewRealFS(base, false),
			},
			[]string{
				".",
				"a",
				"a/aa",
				"a/aa/aaa",
				"a/aa/aab",
				"a/aa/aac",
				"a/ab",
				"a/ab/aba",
				"a/ab/abb",
				"a/ab/abc",
				"a/ac",
				"a/ac/aca",
				"a/ac/acb",
				"a/ac/acc",
			},
			false,
		},
		{
			"scan /, accept only three-triples",
			func(path string, root walker.FS, depth int) (score float64, cont bool, err error) {
				return walker.ScanMatch(depth == 3), true, nil
			},
			walker.Params{
				Root: walker.NewRealFS(base, false),
			},
			[]string{

				"a/aa/aaa",
				"a/aa/aab",
				"a/aa/aac",
				"a/ab/aba",
				"a/ab/abb",
				"a/ab/abc",
				"a/ac/aca",
				"a/ac/acb",
				"a/ac/acc",
			},
			false,
		},
		{
			"scan /, stop inside '/ab'",
			func(pth string, root walker.FS, depth int) (score float64, cont bool, err error) {
				return walker.ScanMatch(true), trimPath(pth) != testutil.ToOSPath("a/ab"), nil
			},
			walker.Params{
				Root: walker.NewRealFS(base, false),
			},
			[]string{
				".",
				"a",
				"a/aa",
				"a/aa/aaa",
				"a/aa/aab",
				"a/aa/aac",
				"a/ab",
				"a/ac",
				"a/ac/aca",
				"a/ac/acb",
				"a/ac/acc",
			},
			false,
		},
		{
			"scan a/aa, don't follow symlinks",
			nil,
			walker.Params{
				Root:       walker.NewRealFS(filepath.Join(base, "a", "aa"), false),
				ExtraRoots: []walker.FS{walker.NewRealFS(filepath.Join(base, "a", "ac"), false)},
			},
			[]string{
				"a/aa",
				"a/aa/aaa",
				"a/aa/aab",
				"a/aa/aac",
				"a/ac",
				"a/ac/aca",
				"a/ac/acb",
				"a/ac/acc",
			},
			false,
		},
		{
			"scan a/aa and extra roots, don't follow links",
			nil,
			walker.Params{
				Root: walker.NewRealFS(filepath.Join(base, "a", "aa"), false),
			},
			[]string{
				"a/aa",
				"a/aa/aaa",
				"a/aa/aab",
				"a/aa/aac",
			},
			false,
		},
		{
			"scan a/aa, follow symlinks", // a/aa/linked links to the root
			nil,
			walker.Params{
				Root: walker.NewRealFS(filepath.Join(base, "a", "aa"), true),
			},
			[]string{
				".",
				"a",
				"a/aa",
				"a/aa/aaa",
				"a/aa/aab",
				"a/aa/aac",
				"a/ab",
				"a/ab/aba",
				"a/ab/abb",
				"a/ab/abc",
				"a/ac",
				"a/ac/aca",
				"a/ac/acb",
				"a/ac/acc",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := walker.Scan(tt.visit, tt.params)
			trimAll(got)
			testutil.ToOSPaths(tt.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanMatch(t *testing.T) {
	t.Parallel()

	type args struct {
		value bool
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"true", args{value: true}, 1},
		{"false", args{value: false}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := walker.ScanMatch(tt.args.value); got != tt.want {
				t.Errorf("ScanMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
