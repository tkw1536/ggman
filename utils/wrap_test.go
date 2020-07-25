package utils

import (
	"reflect"
	"testing"
)

func TestWrapLine(t *testing.T) {
	type args struct {
		line   string
		length int
	}
	tests := []struct {
		name      string
		args      args
		wantLines []string
	}{

		{"wrap short word", args{"prefix hello", 5}, []string{"prefix", "hello"}},
		{"wrap long word", args{"prefix helloworldimtoolong", 5}, []string{"prefix", "helloworldimtoolong"}},

		{"wrap one-word-per-line with normal spaces", args{"hello world beautiful you are", 5}, []string{"hello", "world", "beautiful", "you", "are"}},
		{"wrap one-word-per-line with weird spaces", args{"    hello    world    beautiful you are", 5}, []string{"hello", "world", "beautiful", "you", "are"}},

		{"wrap text normally", args{"hello world beautiful you are", 20}, []string{"hello world", "beautiful you are"}},
		{"wrap text removing spaces", args{"    hello    world    beautiful you are   ", 20}, []string{"hello world", "beautiful you are"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLines := WrapLine(tt.args.line, tt.args.length); !reflect.DeepEqual(gotLines, tt.wantLines) {
				t.Errorf("WrapLine() = %v, want %v", gotLines, tt.wantLines)
			}
		})
	}
}

func TestWrapLinePreserve(t *testing.T) {
	type args struct {
		line   string
		length int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"wrap short word", args{"prefix hello", 5}, []string{"prefix", "hello"}},
		{"wrap long word", args{"prefix helloworldimtoolong", 5}, []string{"prefix", "helloworldimtoolong"}},

		{"wrap one-word-per-line with normal spaces", args{"hello world beautiful you are", 5}, []string{"hello", "world", "beautiful", "you", "are"}},
		{"wrap one-word-per-line with weird spaces", args{"    hello    world    beautiful you are", 5}, []string{"    hello", "    world", "    beautiful", "    you", "    are"}},

		{"wrap text normally", args{"hello world beautiful you are", 20}, []string{"hello world", "beautiful you are"}},
		{"wrap text removing spaces", args{"    hello    world    beautiful you are   ", 20}, []string{"    hello world", "    beautiful you", "    are"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WrapLinePreserve(tt.args.line, tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WrapLinePreserve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapString(t *testing.T) {
	type args struct {
		s      string
		length int
	}
	tests := []struct {
		name      string
		args      args
		wantLines []string
	}{
		{"wrap linux lines", args{"hello \nworld beautiful\nyou are", 5}, []string{"hello", "world", "beautiful", "you", "are"}},
		{"wrap windows lines", args{"hello \r\nworld beautiful\r\nyou are", 5}, []string{"hello", "world", "beautiful", "you", "are"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLines := WrapString(tt.args.s, tt.args.length); !reflect.DeepEqual(gotLines, tt.wantLines) {
				t.Errorf("WrapString() = %v, want %v", gotLines, tt.wantLines)
			}
		})
	}
}

func TestWrapStringPreserve(t *testing.T) {
	type args struct {
		s      string
		length int
	}
	tests := []struct {
		name      string
		args      args
		wantLines []string
	}{
		{"wrap linux lines", args{" hello \n  world beautiful\n   you are", 5}, []string{" hello", "  world", "  beautiful", "   you", "   are"}},
		{"wrap windows lines", args{" hello \r\n  world beautiful\r\n   you are", 5}, []string{" hello", "  world", "  beautiful", "   you", "   are"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLines := WrapStringPreserve(tt.args.s, tt.args.length); !reflect.DeepEqual(gotLines, tt.wantLines) {
				t.Errorf("WrapStringPreserve() = %v, want %v", gotLines, tt.wantLines)
			}
		})
	}
}
