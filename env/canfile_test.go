package env_test

//spellchecker:words reflect strings testing github ggman
import (
	"reflect"
	"strings"
	"testing"

	"github.com/tkw1536/ggman/env"
)

//spellchecker:words canfile

func TestCanLine_UnmarshalText(t *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name    string
		args    args
		wantCl  *env.CanLine
		wantErr bool
	}{
		{"reading pattern-only line", args{[]byte("git@^:$.git")}, &env.CanLine{"", "git@^:$.git"}, false},
		{"reading normal line", args{[]byte("* git@^:$.git")}, &env.CanLine{"*", "git@^:$.git"}, false},
		{"reading line with extra args", args{[]byte("* git@^:$.git extra stuff")}, &env.CanLine{"*", "git@^:$.git"}, false},
		{"empty line is not read", args{[]byte("")}, &env.CanLine{}, true},
		{"comment line is not read", args{[]byte("  //* git@^:$.git extra stuff")}, &env.CanLine{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &env.CanLine{}
			if err := cl.UnmarshalText(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("CanLine.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(cl, tt.wantCl) {
				t.Errorf("CanLine.Unmarshal() = %v, want %v", cl, tt.wantCl)
			}
		})
	}
}

func TestCanFile_ReadFrom(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantCF  env.CanFile
		wantErr bool
	}{
		{
			name:    "empty",
			src:     "",
			wantCF:  env.CanFile(nil),
			wantErr: false,
		},
		{
			name: "canfile with several lines",
			src: `
# for anything on git.example.com, clone with https
^git.example.com https://$.git

# for anything on git2.example.com leave the urls unchanged
^git2.example.com $$

# by default, clone via ssh
git@^:$.git
`,
			wantCF: env.CanFile{
				env.CanLine{Pattern: "^git.example.com", Canonical: "https://$.git"},
				env.CanLine{Pattern: "^git2.example.com", Canonical: "$$"},
				env.CanLine{Pattern: "", Canonical: "git@^:$.git"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cf env.CanFile

			_, gotErr := cf.ReadFrom(strings.NewReader(tt.src))
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("CanFile.ReadFrom() error = %v, wantErr %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(cf, tt.wantCF) {
				t.Errorf("CanLine.ReadFrom() = %#v, want %#v", cf, tt.wantCF)
			}
		})
	}
}
