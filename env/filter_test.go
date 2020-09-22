package env

import "testing"

func TestMatches(t *testing.T) {
	type args struct {
		pattern string
		s       string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// matching the empty pattern
		{"EmptyPattern", args{"", ""}, true},

		// matching one-component parts of a/b/c
		{"oneComponentStart", args{"a", "a/b/c"}, true},
		{"oneComponentMiddle", args{"b", "a/b/c"}, true},
		{"oneComponentEnd", args{"c", "a/b/c"}, true},
		{"oneComponentNot", args{"d", "a/b/c"}, false},

		// matching constant sub-paths
		{"twoComponentsConst", args{"b/c", "a/b/c/d/e/f"}, true},
		{"noTwoComponentsConst", args{"f/g", "a/b/c/d/e/f"}, false},

		// variable sub-paths
		{"variableSubPathPositive", args{"b/*/d", "a/b/c/d/e/f"}, true},
		{"variableSubPathNegative", args{"b/*/c", "a/b/c/d/e/f"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseURL(tt.args.s).Matches(tt.args.pattern); got != tt.want {
				t.Errorf("ParseURL().Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
