//spellchecker:words parseurl
package parseurl_test

//spellchecker:words reflect testing github ggman internal parseurl
import (
	"reflect"
	"testing"

	url "github.com/tkw1536/ggman/internal/parseurl"
)

var slashSplitTests = []struct {
	name  string
	input string
	want  []string
}{

	{
		"empty input",
		"",
		[]string{},
	},
	{
		"only slashes",
		"/////",
		[]string{},
	},

	{
		"regular components",
		"a/b/c",
		[]string{"a", "b", "c"},
	},
	{
		"slash at the start",
		"/a/b/c",
		[]string{"a", "b", "c"}},
	{
		"slash at the end",
		"a/b/c/",
		[]string{"a", "b", "c"},
	},
	{
		"slash at start and end",
		"/a/b/c/",
		[]string{"a", "b", "c"},
	},

	{
		"repeated separator components",
		"a//b//c",
		[]string{"a", "b", "c"},
	},
	{
		"repeated slash at the start",
		"//a/b/c",
		[]string{"a", "b", "c"},
	},
	{
		"repeated slash at the end",
		"a/b/c//",
		[]string{"a", "b", "c"},
	},
	{
		"repeated slash at start and end",
		"//a/b/c//",
		[]string{"a", "b", "c"},
	},
}

func TestSplitNonEmptyRune(t *testing.T) {
	t.Parallel()

	for _, tt := range slashSplitTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := url.SplitNonEmpty(tt.input, '/', nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitNonEmptyRune() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountNonRepeat(t *testing.T) {
	t.Parallel()

	for _, tt := range slashSplitTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wantCount := len(tt.want)

			if gotCount := url.CountNonEmptySplit(tt.input, '/'); gotCount != wantCount {
				t.Errorf("CountNonEmptySplit() = %v, want %v", gotCount, wantCount)
			}
		})
	}
}
