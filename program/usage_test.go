package program

import (
	"reflect"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/program/meta"
)

func TestProgram_MainUsage(t *testing.T) {
	program := makeProgram()

	program.Register(makeEchoCommand("a"))
	program.Register(makeEchoCommand("c"))
	program.Register(makeEchoCommand("b"))

	got := program.MainUsage()
	want := meta.Meta{Executable: "exe", Command: "", Description: "something something dark side", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}, {FieldName: "GlobalOne", Short: []string{"a"}, Long: []string{"global-one"}, Required: false, Value: "", Usage: "", Default: ""}, {FieldName: "GlobalTwo", Short: []string{"b"}, Long: []string{"global-two"}, Required: false, Value: "", Usage: "", Default: ""}}, CommandFlags: []meta.Flag(nil), Positional: meta.Positional{Value: "", Description: "", Min: 0, Max: 0, IncludeUnknown: false}, Commands: []string{"a", "b", "c"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Program.MainUsage() = %#v, want %#v", got, want)
	}
}

func TestProgram_CommandUsage(t *testing.T) {

	program := makeProgram()

	// define requirements to allow only the Global1 (or any) arguments
	reqOne := tRequirements(func(flag meta.Flag) bool {
		return flag.FieldName == "Global1"
	})

	// define requirements to allow anything
	reqAny := tRequirements(func(flag meta.Flag) bool { return true })

	parser := flags.NewParser(&struct {
		Boolean bool `short:"b" value-name:"random" long:"bool" description:"a random boolean argument with short"`
		Int     int  `long:"int" value-name:"dummy" description:"a dummy integer flag" default:"12"`
	}{}, flags.Default)

	type args struct {
		Command     string
		Description string
		Requirement tRequirements
		Positional  meta.Positional
	}
	tests := []struct {
		name string
		args args
		want meta.Meta
	}{
		{
			"command without args and allowing all globals",
			args{Command: "cmd", Requirement: reqAny, Positional: meta.Positional{}},
			meta.Meta{Executable: "exe", Command: "cmd", Description: "", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}, {FieldName: "GlobalOne", Short: []string{"a"}, Long: []string{"global-one"}, Required: false, Value: "", Usage: "", Default: ""}, {FieldName: "GlobalTwo", Short: []string{"b"}, Long: []string{"global-two"}, Required: false, Value: "", Usage: "", Default: ""}}, CommandFlags: []meta.Flag{{FieldName: "Boolean", Short: []string{"b"}, Long: []string{"bool"}, Required: false, Value: "random", Usage: "a random boolean argument with short", Default: ""}, {FieldName: "Int", Short: []string(nil), Long: []string{"int"}, Required: false, Value: "dummy", Usage: "a dummy integer flag", Default: "12"}}, Positional: meta.Positional{Value: "", Description: "", Min: 0, Max: 0, IncludeUnknown: false}, Commands: []string(nil)},
		},

		{
			"command without args and allowing only global1",
			args{Command: "cmd", Requirement: reqOne, Positional: meta.Positional{Description: "usage", Value: "META"}},
			meta.Meta{Executable: "exe", Command: "cmd", Description: "", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}}, CommandFlags: []meta.Flag{{FieldName: "Boolean", Short: []string{"b"}, Long: []string{"bool"}, Required: false, Value: "random", Usage: "a random boolean argument with short", Default: ""}, {FieldName: "Int", Short: []string(nil), Long: []string{"int"}, Required: false, Value: "dummy", Usage: "a dummy integer flag", Default: "12"}}, Positional: meta.Positional{Value: "META", Description: "usage", Min: 0, Max: 0, IncludeUnknown: false}, Commands: []string(nil)},
		},

		{
			"command with max finite args",
			args{Command: "cmd", Requirement: reqOne, Positional: meta.Positional{Max: 4, Description: "usage", Value: "META"}},
			meta.Meta{Executable: "exe", Command: "cmd", Description: "", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}}, CommandFlags: []meta.Flag{{FieldName: "Boolean", Short: []string{"b"}, Long: []string{"bool"}, Required: false, Value: "random", Usage: "a random boolean argument with short", Default: ""}, {FieldName: "Int", Short: []string(nil), Long: []string{"int"}, Required: false, Value: "dummy", Usage: "a dummy integer flag", Default: "12"}}, Positional: meta.Positional{Value: "META", Description: "usage", Min: 0, Max: 4, IncludeUnknown: false}, Commands: []string(nil)},
		},

		{
			"command with finite args",
			args{Command: "cmd", Requirement: reqOne, Positional: meta.Positional{Min: 1, Max: 2, Description: "usage", Value: "META"}},
			meta.Meta{Executable: "exe", Command: "cmd", Description: "", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}}, CommandFlags: []meta.Flag{{FieldName: "Boolean", Short: []string{"b"}, Long: []string{"bool"}, Required: false, Value: "random", Usage: "a random boolean argument with short", Default: ""}, {FieldName: "Int", Short: []string(nil), Long: []string{"int"}, Required: false, Value: "dummy", Usage: "a dummy integer flag", Default: "12"}}, Positional: meta.Positional{Value: "META", Description: "usage", Min: 1, Max: 2, IncludeUnknown: false}, Commands: []string(nil)},
		},

		{
			"command with infinite args",
			args{Command: "cmd", Requirement: reqOne, Positional: meta.Positional{Min: 1, Max: -1, Description: "usage", Value: "META"}},
			meta.Meta{Executable: "exe", Command: "cmd", Description: "", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}}, CommandFlags: []meta.Flag{{FieldName: "Boolean", Short: []string{"b"}, Long: []string{"bool"}, Required: false, Value: "random", Usage: "a random boolean argument with short", Default: ""}, {FieldName: "Int", Short: []string(nil), Long: []string{"int"}, Required: false, Value: "dummy", Usage: "a dummy integer flag", Default: "12"}}, Positional: meta.Positional{Value: "META", Description: "usage", Min: 1, Max: -1, IncludeUnknown: false}, Commands: []string(nil)},
		},

		{
			"command with description",
			args{Command: "cmd", Description: "A fake command", Requirement: reqOne, Positional: meta.Positional{Min: 1, Max: -1, Description: "usage", Value: "META"}},
			meta.Meta{Executable: "exe", Command: "cmd", Description: "A fake command", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}}, CommandFlags: []meta.Flag{{FieldName: "Boolean", Short: []string{"b"}, Long: []string{"bool"}, Required: false, Value: "random", Usage: "a random boolean argument with short", Default: ""}, {FieldName: "Int", Short: []string(nil), Long: []string{"int"}, Required: false, Value: "dummy", Usage: "a dummy integer flag", Default: "12"}}, Positional: meta.Positional{Value: "META", Description: "usage", Min: 1, Max: -1, IncludeUnknown: false}, Commands: []string(nil)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := iContext{
				Args: iArguments{
					Command: tt.args.Command,
				},

				commandParser: parser,

				Description: iDescription{
					Command:      tt.args.Command,
					Description:  tt.args.Description,
					Positional:   tt.args.Positional,
					Requirements: tt.args.Requirement,
				},
			}
			got := program.CommandUsage(context)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Program.CommandUsage() = %#v, want %v", got, tt.want)
			}
		})
	}
}
