package meta

import (
	"reflect"
	"strings"
	"testing"

	"github.com/jessevdk/go-flags"
)

func TestNewFlag(t *testing.T) {
	tests := []struct {
		name string
		opt  *flags.Option
		want Flag
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
			Flag{
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
			Flag{
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

			got := NewFlag(tt.opt)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFlag() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestFlag_WriteSpecTo(t *testing.T) {
	tests := []struct {
		name string
		flag Flag
		want string
	}{
		{
			"long only optional option",
			Flag{Long: []string{"long"}},
			"[--long]",
		},
		{
			"short and long optional option",
			Flag{Short: []string{"s"}, Long: []string{"long"}},
			"[--long|-s]",
		},
		{
			"short and long named optional option",
			Flag{Value: "name", Short: []string{"s"}, Long: []string{"long"}},
			"[--long|-s name]",
		},

		{
			"long only required option",
			Flag{Long: []string{"long"}, Required: true},
			"--long",
		},
		{
			"short and long required option",
			Flag{Short: []string{"s"}, Long: []string{"long"}, Required: true},
			"--long|-s",
		},
		{
			"short and long named required option",
			Flag{Value: "name", Short: []string{"s"}, Long: []string{"long"}, Required: true},
			"--long|-s name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			tt.flag.WriteSpecTo(&builder)
			if got := builder.String(); got != tt.want {
				t.Errorf("Flag.WriteSpecTo() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFlag_WriteLongSpecTo(t *testing.T) {
	tests := []struct {
		name string
		flag Flag
		want string
	}{
		{
			"long only option",
			Flag{Long: []string{"long"}},
			"--long",
		},
		{
			"short and long option",
			Flag{Short: []string{"s"}, Long: []string{"long"}},
			"-s, --long",
		},
		{
			"short and long named option",
			Flag{Value: "name", Short: []string{"s"}, Long: []string{"long"}},
			"-s, --long name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			tt.flag.WriteLongSpecTo(&builder)
			if got := builder.String(); got != tt.want {
				t.Errorf("Flag.WriteLongSpecTo() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFlag_WriteMessageTo(t *testing.T) {
	tests := []struct {
		name string
		flag Flag
		want string
	}{
		{
			"long only option",
			Flag{Usage: "a long option", Long: []string{"long"}},
			"\n\n   --long\n      a long option",
		},
		{
			"short and long option",
			Flag{Usage: "a long or short option", Short: []string{"s"}, Long: []string{"long"}},
			"\n\n   -s, --long\n      a long or short option",
		},
		{
			"short and long named option",
			Flag{Usage: "this one is named", Value: "name", Short: []string{"s"}, Long: []string{"long"}},
			"\n\n   -s, --long name\n      this one is named",
		},
		{
			"short and long named option with default",
			Flag{Usage: "this one is named", Value: "name", Short: []string{"s"}, Long: []string{"long"}, Default: "default"},
			"\n\n   -s, --long name\n      this one is named (default default)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			tt.flag.WriteMessageTo(&builder)
			if got := builder.String(); got != tt.want {
				t.Errorf("Flag.WriteMessageTo() = %q, want %q", got, tt.want)
			}
		})
	}
}
