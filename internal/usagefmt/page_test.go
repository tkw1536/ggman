package usagefmt

import "testing"

func TestPage_String(t *testing.T) {
	type fields struct {
		MainName    string
		MainOpts    []Opt
		Description string
		SubCommands []string
		SubName     string
		SubOpts     []Opt
		MetaName    string
		MetaMin     int
		MetaMax     int
		Usage       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"main executable page",
			fields{
				MainName: "cmd",
				MainOpts: []Opt{
					FakeOpt{
						required: true,

						short: []string{"g"},
						long:  []string{"global"},

						value: "name",
						usage: "A global argument",
						dflt:  "",
					},
					FakeOpt{
						required: false,

						short: []string{"q"},
						long:  []string{"quiet"},
						usage: "Be quiet",
						dflt:  "false",
					},
				},
				Description: "Do something interesting",
				SubCommands: []string{"a", "b", "c"},
			},
			"Usage: cmd --global|-g name [--quiet|-q] [--] COMMAND [ARGS...]\n\nDo something interesting\n\n   -g, --global name\n      A global argument\n\n   -q, --quiet\n      Be quiet (default false)\n\n   COMMAND [ARGS...]\n      Command to call. One of \"a\", \"b\", \"c\". See individual commands for more help.",
		},
		{
			"sub executable page",
			fields{
				MainName: "cmd",
				MainOpts: []Opt{
					FakeOpt{
						required: true,

						short: []string{"g"},
						long:  []string{"global"},

						value: "name",
						usage: "A global argument",
						dflt:  "",
					},
					FakeOpt{
						required: false,

						short: []string{"q"},
						long:  []string{"quiet"},
						usage: "Be quiet",
						dflt:  "false",
					},
				},
				Description: "Do something local",
				SubName:     "sub",
				SubOpts: []Opt{
					FakeOpt{
						required: true,

						short: []string{"d"},
						long:  []string{"dud"},

						value: "dud",
						usage: "A local argument",
						dflt:  "",
					},
					FakeOpt{
						required: false,

						short: []string{"s"},
						long:  []string{"silent"},
						usage: "Be silent",
						dflt:  "true",
					},
				},
				MetaName: "op",
				MetaMin:  1,
				MetaMax:  -1,
				Usage:    "Operations to make",
			},
			"Usage: cmd --global|-g name [--quiet|-q] [--] sub --dud|-d dud [--silent|-s] [--] op [op ...]\n\nDo something local\n\nGlobal Arguments:\n\n   -g, --global name\n      A global argument\n\n   -q, --quiet\n      Be quiet (default false)\n\nCommand Arguments:\n\n   -d, --dud dud\n      A local argument\n\n   -s, --silent\n      Be silent (default true)\n\n   op [op ...]\n      Operations to make",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := Page{
				MainName:    tt.fields.MainName,
				MainOpts:    tt.fields.MainOpts,
				Description: tt.fields.Description,
				SubCommands: tt.fields.SubCommands,
				SubName:     tt.fields.SubName,
				SubOpts:     tt.fields.SubOpts,
				MetaName:    tt.fields.MetaName,
				MetaMin:     tt.fields.MetaMin,
				MetaMax:     tt.fields.MetaMax,
				Usage:       tt.fields.Usage,
			}
			if got := page.String(); got != tt.want {
				t.Errorf("Page.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
