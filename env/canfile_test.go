package env

//spellchecker:words reflect testing
import (
	"reflect"
	"testing"
)

func TestCanLine_Unmarshal(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantCl  *CanLine
		wantErr bool
	}{
		{"reading pattern-only line", args{"git@^:$.git"}, &CanLine{"", "git@^:$.git"}, false},
		{"reading normal line", args{"* git@^:$.git"}, &CanLine{"*", "git@^:$.git"}, false},
		{"reading line with extra args", args{"* git@^:$.git extra stuff"}, &CanLine{"*", "git@^:$.git"}, false},
		{"empty line is not read", args{""}, &CanLine{}, true},
		{"comment line is not read", args{"  //* git@^:$.git extra stuff"}, &CanLine{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &CanLine{}
			if err := cl.Unmarshal(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("CanLine.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(cl, tt.wantCl) {
				t.Errorf("CanLine.Unmarshal() = %v, want %v", cl, tt.wantCl)
			}
		})
	}
}
