package walker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/ggman/program/lib/testlib"
)

func Test_IsDirectory(t *testing.T) {
	dir := testlib.TempDirAbs(t)

	// make a symlink to a directory
	dirlink := filepath.Join(dir, "dirlink")
	if err := os.Symlink(dir, dirlink); err != nil {
		panic(err)
	}

	// make a file
	file := filepath.Join(dir, "file")
	if err := os.WriteFile(file, nil, os.ModePerm); err != nil {
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
