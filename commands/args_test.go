package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseArgs tests the ParseArgs Function
func TestParseArgs(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		input    []string
		expected *GGArgs
	}{
		// no arguments => parsing fails
		{[]string{}, nil},

		// command without arguments => ok
		{[]string{"cmd"}, &GGArgs{"cmd", "", false, []string{}}},

		// command with help => help flag set
		{[]string{"help", "cmd"}, &GGArgs{"", "", true, []string{"cmd"}}},
		{[]string{"--help", "cmd"}, &GGArgs{"", "", true, []string{"cmd"}}},
		{[]string{"-h", "cmd"}, &GGArgs{"", "", true, []string{"cmd"}}},

		// command with arguments => ok
		{[]string{"cmd", "a1", "a2"}, &GGArgs{"cmd", "", false, []string{"a1", "a2"}}},

		// command with help => passed through
		{[]string{"cmd", "help", "a1"}, &GGArgs{"cmd", "", false, []string{"help", "a1"}}},
		{[]string{"cmd", "--help", "a1"}, &GGArgs{"cmd", "", false, []string{"--help", "a1"}}},
		{[]string{"cmd", "-h", "a1"}, &GGArgs{"cmd", "", false, []string{"-h", "a1"}}},

		// only a for => parsing fails
		{[]string{"for"}, nil},
		{[]string{"--for"}, nil},
		{[]string{"-f"}, nil},

		// for without command => parsing fails
		{[]string{"for", "match"}, nil},
		{[]string{"--for", "match"}, nil},
		{[]string{"-f", "match"}, nil},

		// for with command => ok
		{[]string{"for", "match", "cmd"}, &GGArgs{"cmd", "match", false, []string{}}},
		{[]string{"--for", "match", "cmd"}, &GGArgs{"cmd", "match", false, []string{}}},
		{[]string{"-f", "match", "cmd"}, &GGArgs{"cmd", "match", false, []string{}}},

		// for with command and arguments => ok
		{[]string{"for", "match", "cmd", "a1", "a2"}, &GGArgs{"cmd", "match", false, []string{"a1", "a2"}}},
		{[]string{"--for", "match", "cmd", "a1", "a2"}, &GGArgs{"cmd", "match", false, []string{"a1", "a2"}}},
		{[]string{"-f", "match", "cmd", "a1", "a2"}, &GGArgs{"cmd", "match", false, []string{"a1", "a2"}}},
	}

	for _, row := range table {
		res, err := ParseArgs(row.input)
		if row.expected == nil {
			assert.Nil(res)
			assert.NotEqual("", err)
		} else {
			assert.Equal(row.expected, res)
			assert.Equal("", err)
		}
	}
}

func TestGGArgs_ParseSingleFlag(t *testing.T) {
	type fields struct {
		Command string
		Pattern string
		Args    []string
	}
	type args struct {
		flag string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue bool
		wantErr   bool
	}{
		// giving no arguments
		{"no arguments given", fields{"cmd", "", []string{}}, args{"--test"}, false, false},
		{"right argument given", fields{"cmd", "", []string{"--test"}}, args{"--test"}, true, false},
		{"wrong argument given", fields{"cmd", "", []string{"--fake"}}, args{"--test"}, false, true},
		{"too many arguments", fields{"cmd", "", []string{"--fake", "--untrue"}}, args{"--test"}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := &GGArgs{
				Command: tt.fields.Command,
				Pattern: tt.fields.Pattern,
				Args:    tt.fields.Args,
			}
			gotValue, gotErr := parsed.ParseSingleFlag(tt.args.flag)
			if gotValue != tt.wantValue {
				t.Errorf("GGArgs.ParseSingleFlag() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotErr != tt.wantErr {
				t.Errorf("GGArgs.ParseSingleFlag() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
