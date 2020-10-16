package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/tkw1536/ggman/testutil"
)

func TestScan(t *testing.T) {

	base, cleanup := testutil.TempDir()
	defer cleanup()

	// setup a directory structure for testing.
	// Make mkdir and symlink utility methods for this.

	mkdir := func(s string) {
		err := os.MkdirAll(filepath.Join(base, s), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	symlink := func(oldname, newname string) {
		err := os.Symlink(filepath.Join(base, oldname), filepath.Join(base, newname))
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
		options ScanOptions
		want    []string
		wantErr bool
	}{
		{
			"scan /",

			ScanOptions{
				Root:        base,
				FollowLinks: false,
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

			ScanOptions{
				Root:        base,
				FollowLinks: false,
				Filter: func(path string) (match, cont bool) {
					return strings.Count(trimPath(path), string(filepath.Separator)) == 2, true
				},
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

			ScanOptions{
				Root:        base,
				FollowLinks: false,
				Filter: func(path string) (match, cont bool) {
					return true, trimPath(path) != ToOSPath("a/ab")
				},
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

			ScanOptions{
				Root:        filepath.Join(base, "a", "aa"),
				FollowLinks: false,
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

			ScanOptions{
				Root:        filepath.Join(base, "a", "aa"),
				FollowLinks: true,
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
			got, err := Scan(tt.options)
			trimAll(got)
			ToOSPaths(tt.want)
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

func Test_IsDirectory(t *testing.T) {
	dir, cleanup := testutil.TempDir()
	defer cleanup()

	// make a symlink to a directory
	dirlink := filepath.Join(dir, "dirlink")
	if err := os.Symlink(dir, dirlink); err != nil {
		panic(err)
	}

	// make a file
	file := filepath.Join(dir, "file")
	if err := ioutil.WriteFile(file, nil, os.ModePerm); err != nil {
		panic(err)
	}

	// make a symlink to a file
	filelink := filepath.Join(dir, "filelink")
	if err := os.Symlink(file, filelink); err != nil {
		panic(err)
	}

	type args struct {
		path       string
		allowLinks bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"directory without allowLinks", args{dir, false}, true, false},
		{"directory with allowLinks", args{dir, true}, true, false},

		{"directory link without allowLinks", args{dirlink, false}, false, false},
		{"directory link with allowLinks", args{dirlink, true}, true, false},

		{"file without allowLinks", args{file, false}, false, false},
		{"file with allowLinks", args{file, true}, false, false},

		{"file link without allowLinks", args{filelink, false}, false, false},
		{"file link with allowLinks", args{filelink, true}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsDirectory(tt.args.path, tt.args.allowLinks)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}
