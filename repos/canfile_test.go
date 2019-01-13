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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCl := ReadCanLine(tt.args.line); !reflect.DeepEqual(gotCl, tt.wantCl) {
				t.Errorf("ReadCanLine() = %v, want %v", gotCl, tt.wantCl)
			}
		})
	}
}
