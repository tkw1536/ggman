package program

import (
	"reflect"
	"testing"

	"github.com/jessevdk/go-flags"
)

func Test_parseFlagNames(t *testing.T) {
	type args struct {
		err *flags.Error
	}
	tests := []struct {
		name      string
		args      args
		wantNames []string
		wantOk    bool
	}{
		{
			"unix flag message",
			args{
				&flags.Error{
					Message: "expected argument for flag `-f, --for'",
				},
			},
			[]string{"f", "for"},
			true,
		},
		{
			"windows flag message",
			args{
				&flags.Error{
					Message: "expected argument for flag `/f, /for'",
				},
			},
			[]string{"f", "for"},
			true,
		},

		{
			"no message",
			args{
				&flags.Error{
					Message: "",
				},
			},
			nil,
			false,
		},

		{
			"only beginning",
			args{
				&flags.Error{
					Message: "expected argument for flag `-f, --for",
				},
			},
			nil,
			false,
		},

		{
			"only end",
			args{
				&flags.Error{
					Message: "expected argument for flag -f, --for'",
				},
			},
			nil,
			false,
		},

		{
			"wrong order",
			args{
				&flags.Error{
					Message: "expected argument for flag '-f, --for`",
				},
			},
			nil,
			false,
		},

		{
			"empty flags",
			args{
				&flags.Error{
					Message: "expected argument for flag `'",
				},
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNames, gotOk := parseFlagNames(tt.args.err)
			if !reflect.DeepEqual(gotNames, tt.wantNames) {
				t.Errorf("parseFlagNames() gotNames = %v, want %v", gotNames, tt.wantNames)
			}
			if gotOk != tt.wantOk {
				t.Errorf("parseFlagNames() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
