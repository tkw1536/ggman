package usagefmt

import (
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/internal/text"
)

// FakeOpt implements Opt for testing purposes
type FakeOpt struct {
	required           bool
	short, long        []string
	value, usage, dflt string
}

var _ Opt = (*FakeOpt)(nil)

func (f FakeOpt) Required() bool  { return f.required }
func (f FakeOpt) Short() []string { return f.short }
func (f FakeOpt) Long() []string  { return f.long }
func (f FakeOpt) Value() string   { return f.value }
func (f FakeOpt) Usage() string   { return f.usage }
func (f FakeOpt) Default() string { return f.dflt }

func TestNewOpt(t *testing.T) {
	tests := []struct {
		name string
		opt  *flags.Option
		want FakeOpt
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
			FakeOpt{
				required: true,

				short: []string{"s"},
				long:  []string{"long"},

				value: "test",
				usage: "something",
				dflt:  "",
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
			FakeOpt{
				required: false,

				long: []string{"long"},

				value: "test",
				usage: "something",
				dflt:  "a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			opt := NewOpt(tt.opt)

			if tt.want.Required() != opt.Required() {
				t.Errorf("NewOpt().Required() = %v, want %v", tt.want.Required(), opt.Required())
			}

			if !text.SliceEquals(tt.want.Short(), opt.Short()) {
				t.Errorf("NewOpt().Short() = %v, want %v", tt.want.Short(), opt.Short())
			}

			if !text.SliceEquals(tt.want.Long(), opt.Long()) {
				t.Errorf("NewOpt().Long() = %v, want %v", tt.want.Long(), opt.Long())
			}

			if tt.want.Value() != opt.Value() {
				t.Errorf("NewOpt().Required() = %v, want %v", tt.want.Value(), opt.Value())
			}

			if tt.want.Usage() != opt.Usage() {
				t.Errorf("NewOpt().Usage() = %v, want %v", tt.want.Usage(), opt.Usage())
			}

			if tt.want.Default() != opt.Default() {
				t.Errorf("NewOpt().Default() = %v, want %v", tt.want.Default(), opt.Default())
			}
		})
	}
}
