package path

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/pkglib/testlib"
)

func TestSameFile(t *testing.T) {
	touch := func(path string) {
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	symlink := func(old, new string) {
		if err := os.MkdirAll(filepath.Dir(new), fs.ModeDir&fs.ModePerm); err != nil {
			panic(err)
		}
		if err := os.Symlink(old, new); err != nil {
			panic(err)
		}
	}

	d1 := testlib.TempDirAbs(t)

	f1 := filepath.Join(d1, "f1")
	touch(f1)

	alsoF1 := filepath.Join(d1, "same")
	symlink(f1, alsoF1)

	d2 := testlib.TempDirAbs(t)
	alsoD2 := filepath.Join(d2, "nested")
	symlink(d2, alsoD2)

	f2 := filepath.Join(d2, "f2")
	f3 := filepath.Join(d2, "f3")

	d3 := testlib.TempDirAbs(t)
	_ = d3

	tests := []struct {
		name  string
		path1 string
		path2 string
		want  bool
	}{

		{"identical existing files", f1, f1, true},
		{"identical linked files (1)", f1, alsoF1, true},
		{"identical linked files (2)", alsoF1, f1, true},
		{"identical linked files (3)", alsoF1, alsoF1, true},

		{"identical non-existing files (1)", f2, f2, true},
		{"non-identical non-existing files (1)", alsoF1, f2, false},
		{"non-identical non-existing files (2)", f2, f3, false},
		{"non-identical partial existing files", f2, f1, false},
		{"non-identical partial existing files", f1, f2, false},

		{"identical existing directories", d1, d1, true},
		{"identical linked directories", d2, alsoD2, true},
		{"non-identical existing direcories (1)", d2, d1, false},
		{"non-identical existing direcories (2)", d1, d2, false},
		{"non-identical existing direcories (3)", d1, alsoD2, false},
		{"non-identical existing direcories (4)", alsoD2, d1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SameFile(tt.path1, tt.path2); got != tt.want {
				t.Errorf("SameFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
