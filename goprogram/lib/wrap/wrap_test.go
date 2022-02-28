package wrap

import (
	"io"
	"strings"
	"testing"
)

func TestWrapper_WrapLine(t *testing.T) {
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

	var builder strings.Builder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder.Reset()
			wrapper := &Wrapper{Writer: &builder, Length: tt.args.length}

			wrapper.WriteLine(tt.args.line)
			gotLine := builder.String()
			wantLine := strings.Join(tt.wantLines, "\n")

			if gotLine != wantLine {
				t.Errorf("Wrapper.WriteLine() = %q, want %q", gotLine, wantLine)
			}
		})
	}
}

func BenchmarkWrapper_WriteLine(b *testing.B) {
	wrapper := &Wrapper{
		Writer: io.Discard,
		Length: 20,
	}

	for n := 0; n < b.N; n++ {
		wrapper.WriteLine(` Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam eget tortor massa. Nullam gravida massa id dui placerat condimentum. Proin volutpat massa eu enim luctus convallis. Integer a nulla facilisis, convallis elit id, tristique nisi. Duis enim diam, viverra sed quam quis, scelerisque aliquam mi. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Sed consectetur cursus libero, non lobortis mi tempor sit amet. Nullam sapien sapien, imperdiet id cursus non, consequat sed neque. Fusce sollicitudin tortor pulvinar, placerat urna sit amet, luctus tellus. Vivamus sit amet ligula purus. `)
	}
}

func TestWrapper_WriteIndent(t *testing.T) {
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
		{"wrap only spaces", args{"               ", 20}, []string{"               "}},
	}

	var builder strings.Builder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder.Reset()

			wrapper := &Wrapper{
				Writer: &builder,
				Length: tt.args.length,
			}

			wrapper.WriteIndent(tt.args.line)
			gotLine := builder.String()
			wantLine := strings.Join(tt.want, "\n")

			if gotLine != wantLine {
				t.Errorf("Wrapper.WriteIndent() = %q, want %q", gotLine, wantLine)
			}
		})
	}
}

func BenchmarkWrapper_WriteIndent(b *testing.B) {
	wrapper := &Wrapper{
		Writer: io.Discard,
		Length: 20,
	}
	for n := 0; n < b.N; n++ {
		wrapper.WriteIndent(`         Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam eget tortor massa. Nullam gravida massa id dui placerat condimentum. Proin volutpat massa eu enim luctus convallis. Integer a nulla facilisis, convallis elit id, tristique nisi. Duis enim diam, viverra sed quam quis, scelerisque aliquam mi. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Sed consectetur cursus libero, non lobortis mi tempor sit amet. Nullam sapien sapien, imperdiet id cursus non, consequat sed neque. Fusce sollicitudin tortor pulvinar, placerat urna sit amet, luctus tellus. Vivamus sit amet ligula purus. `)
	}
}

func TestWrapper_WriteString(t *testing.T) {
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

	var builder strings.Builder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder.Reset()

			wrapper := &Wrapper{
				Writer: &builder,
				Length: tt.args.length,
			}
			wrapper.WriteString(tt.args.s)
			gotLine := builder.String()
			wantLine := strings.Join(tt.wantLines, "\n")

			if gotLine != wantLine {
				t.Errorf("Wrapper.WriteString() = %q, want %q", gotLine, wantLine)
			}
		})
	}
}
