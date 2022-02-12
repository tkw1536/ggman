package usagefmt

import (
	"reflect"
	"testing"

	"github.com/jessevdk/go-flags"
)

func TestNewOpt(t *testing.T) {
	tests := []struct {
		name string
		opt  *flags.Option
		want Opt
	}{
		{
			"simple option without default",
			&flags.Option{
				Required: true,

				ShortName: 's',
				LongName:  "long",

				ValueName:   "test",
				Description: "something",
				Default:     nil,
			},
			Opt{
				Required: true,

				Short:     []string{"s"},
				Long:      []string{"long"},
				FieldName: "",

				Value:   "test",
				Usage:   "something",
				Default: "",
			},
		},

		{
			"simple option with default",
			&flags.Option{
				Required: false,

				LongName: "long",

				ValueName:   "test",
				Description: "something",
				Default:     []string{"a"},
			},
			Opt{
				Required: false,

				Long: []string{"long"},

				Value:   "test",
				Usage:   "something",
				Default: "a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := NewOpt(tt.opt)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOpt() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
