package repos

import (
	"reflect"
	"testing"
)

func TestReadCanLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		args   args
		wantCl *CanLine
	}{
		{"reading pattern-only line", args{"git@^:$.git"}, &CanLine{"", "git@^:$.git"}},
		{"reading normal line", args{"* git@^:$.git"}, &CanLine{"*", "git@^:$.git"}},
		{"reading line with extra args", args{"* git@^:$.git extra stuff"}, &CanLine{"*", "git@^:$.git"}},
		{"empty line is not read", args{""}, nil},
		{"comment line is not read", args{"  //* git@^:$.git extra stuff"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCl := ReadCanLine(tt.args.line); !reflect.DeepEqual(gotCl, tt.wantCl) {
				t.Errorf("ReadCanLine() = %v, want %v", gotCl, tt.wantCl)
			}
		})
	}
}
