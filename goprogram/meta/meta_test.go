package meta

import (
	"strings"
	"testing"
)

func TestUsage_WriteMessageTo(t *testing.T) {
	tests := []struct {
		name string
		meta Meta
		want string
	}{
		{
			"main executable page",
			Meta{
				Executable:  "cmd",
				Description: "Do something interesting",

				GlobalFlags: []Flag{
					{
						Required: true,

						Short: []string{"g"},
						Long:  []string{"global"},

						Value:   "name",
						Usage:   "A global argument",
						Default: "",
					},
					{
						Required: false,

						Short:   []string{"q"},
						Long:    []string{"quiet"},
						Usage:   "Be quiet",
						Default: "false",
					},
				},
				Commands: []string{"a", "b", "c"},
			},
			"Usage: cmd --global|-g name [--quiet|-q] [--] COMMAND [ARGS...]\n\nDo something interesting\n\n   -g, --global name\n      A global argument\n\n   -q, --quiet\n      Be quiet (default false)\n\n   COMMAND [ARGS...]\n      Command to call. One of \"a\", \"b\", \"c\". See individual commands for more help.",
		},
		{
			"sub executable page",
			Meta{
				Executable:  "cmd",
				Command:     "sub",
				Description: "Do something local",

				GlobalFlags: []Flag{
					{
						Required: true,

						Short: []string{"g"},
						Long:  []string{"global"},

						Value:   "name",
						Usage:   "A global argument",
						Default: "",
					},
					{
						Required: false,

						Short:   []string{"q"},
						Long:    []string{"quiet"},
						Usage:   "Be quiet",
						Default: "false",
					},
				},
				CommandFlags: []Flag{
					{
						Required: true,

						Short: []string{"d"},
						Long:  []string{"dud"},

						Value:   "dud",
						Usage:   "A local argument",
						Default: "",
					},
					{
						Required: false,

						Short:   []string{"s"},
						Long:    []string{"silent"},
						Usage:   "Be silent",
						Default: "true",
					},
				},
				Positional: Positional{
					Value:       "op",
					Description: "Operations to make",
					Min:         1,
					Max:         -1,
				},
			},
			"Usage: cmd --global|-g name [--quiet|-q] [--] sub --dud|-d dud [--silent|-s] [--] op [op ...]\n\nDo something local\n\nGlobal Arguments:\n\n   -g, --global name\n      A global argument\n\n   -q, --quiet\n      Be quiet (default false)\n\nCommand Arguments:\n\n   -d, --dud dud\n      A local argument\n\n   -s, --silent\n      Be silent (default true)\n\n   op [op ...]\n      Operations to make",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			tt.meta.WriteMessageTo(&builder)

			if got := builder.String(); got != tt.want {
				t.Errorf("Usage.WriteMessageTo() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMeta_writeCommandsTo(t *testing.T) {
	type args struct {
		commands []string
	}
	tests := []struct {
		name string
		meta Meta
		want string
	}{
		{"no commands", Meta{Commands: nil}, ""},
		{"single command", Meta{Commands: []string{"a"}}, `"a"`},
		{"multiple commands", Meta{Commands: []string{"a", "b", "c"}}, `"a", "b", "c"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			tt.meta.writeCommandsTo(&builder)
			if got := builder.String(); got != tt.want {
				t.Errorf("Meta.writeCommandsTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
