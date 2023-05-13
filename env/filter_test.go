package env

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/internal/path"
	"github.com/tkw1536/ggman/internal/pattern"
	"github.com/tkw1536/ggman/internal/testutil"
	"github.com/tkw1536/pkglib/testlib"
)

func setupFilterTest(t *testing.T) (root, exampleClonePath, otherClonePath string) {
	root = testlib.TempDirAbs(t)

	exampleClonePath = filepath.Join(root, "example")
	if testutil.NewTestRepoAt(exampleClonePath, "") == nil {
		panic("failed to create test repo")
	}

	otherClonePath = filepath.Join(root, "other")
	if testutil.NewTestRepoAt(otherClonePath, "") == nil {
		panic("failed to create test repo")
	}

	return root, exampleClonePath, otherClonePath
}

func Test_emptyFilter_Score(t *testing.T) {
	root, exampleClonePath, otherClonePath := setupFilterTest(t)

	type args struct {
		env       Env
		clonePath string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"empty filter matches clone path",
			args{env: Env{Root: root}, clonePath: exampleClonePath},
			1,
		},
		{
			"empty filter matches other clone path",
			args{env: Env{Root: root}, clonePath: otherClonePath},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := emptyFilter{}
			if got := e.Score(tt.args.env, tt.args.clonePath); got != tt.want {
				t.Errorf("emptyFilter.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testFilter struct{}

func (testFilter) Score(env Env, clonePath string) float64 { panic("never reached") }

type testFilterWithCandidates struct {
	testFilter
}

func (testFilterWithCandidates) Candidates() []string { return []string{"a", "b", "c"} }

func TestCandidates(t *testing.T) {
	tests := []struct {
		name   string
		filter Filter
		want   []string
	}{
		{
			"candidates of non-candidate-filter is nil",
			testFilter{},
			nil,
		},
		{
			"candidates of candidate-filter calls candidates",
			testFilterWithCandidates{},
			[]string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Candidates(tt.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Candidates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathFilter_Score(t *testing.T) {
	root, exampleClonePath, otherClonePath := setupFilterTest(t)

	type fields struct {
		Paths []string
	}
	type args struct {
		env       Env
		clonePath string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			"non-listed path doesn't match",
			fields{
				Paths: []string{exampleClonePath, otherClonePath},
			},
			args{
				Env{Root: root},
				root,
			},
			-1,
		},
		{
			"non-listed path doesn't match",
			fields{
				Paths: []string{exampleClonePath, otherClonePath},
			},
			args{
				Env{Root: root},
				"/outside/",
			},
			-1,
		},
		{
			"listed path matches (1)",
			fields{
				Paths: []string{exampleClonePath, otherClonePath},
			},
			args{
				Env{Root: root},
				exampleClonePath,
			},
			1,
		},
		{
			"listed path matches (2)",
			fields{
				Paths: []string{exampleClonePath, otherClonePath},
			},
			args{
				Env{Root: root},
				otherClonePath,
			},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := PathFilter{
				Paths: tt.fields.Paths,
			}
			if got := pf.Score(tt.args.env, tt.args.clonePath); got != tt.want {
				t.Errorf("PathFilter.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathFilter_Candidates(t *testing.T) {
	_, exampleClonePath, otherClonePath := setupFilterTest(t)

	type fields struct {
		Paths []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"candidates with list",
			fields{
				Paths: []string{exampleClonePath, otherClonePath},
			},
			[]string{exampleClonePath, otherClonePath},
		},
		{
			"candidates with nil list",
			fields{
				Paths: nil,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := PathFilter{
				Paths: tt.fields.Paths,
			}
			if got := pf.Candidates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PathFilter.Candidates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPatternFilter(t *testing.T) {
	type args struct {
		value string
		fuzzy bool
	}
	tests := []struct {
		name    string
		args    args
		wantPat PatternFilter
	}{
		{
			"a/b (non-fuzzy)",
			args{"a/b", false},
			PatternFilter{
				value:   "a/b",
				pattern: pattern.NewSplitGlobPattern("a/b", ComponentsOf, false),
			},
		},
		{
			"'' (non-fuzzy)",
			args{"", false},
			PatternFilter{
				value:   "",
				pattern: pattern.NewSplitGlobPattern("", ComponentsOf, false),
			},
		},

		{
			"a/b (fuzzy)",
			args{"a/b", true},
			PatternFilter{
				value:   "a/b",
				fuzzy:   true,
				pattern: pattern.NewSplitGlobPattern("a/b", ComponentsOf, true),
			},
		},
		{
			"'' (fuzzy)",
			args{"", true},
			PatternFilter{
				value:   "",
				fuzzy:   true,
				pattern: pattern.NewSplitGlobPattern("", ComponentsOf, true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPat := NewPatternFilter(tt.args.value, tt.args.fuzzy)

			// .Split cannot be compared with reflect
			gotPat.pattern.Split = nil
			tt.wantPat.pattern.Split = nil

			if !reflect.DeepEqual(gotPat, tt.wantPat) {
				t.Errorf("NewPatternFilter() = %v, want %v", gotPat, tt.wantPat)
			}
		})
	}
}

func TestPatternFilter_String(t *testing.T) {
	type fields struct {
		value   string
		pattern pattern.SplitPattern
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"empty pattern",
			fields{
				value:   "",
				pattern: pattern.SplitPattern{},
			},
			"",
		},
		{
			"a/b pattern",
			fields{
				value:   "a/b",
				pattern: pattern.SplitPattern{},
			},
			"a/b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := PatternFilter{
				value:   tt.fields.value,
				pattern: tt.fields.pattern,
			}
			if got := pat.String(); got != tt.want {
				t.Errorf("PatternFilter.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatternFilter_Score(t *testing.T) {
	root := testlib.TempDirAbs(t)

	abc := filepath.Join(root, "a", "b", "c")
	abcdef := filepath.Join(root, "a", "b", "c", "d", "e", "f")

	if testutil.NewTestRepoAt(abc, "/a/b/c") == nil {
		panic("NewTestRepoAt() returned nil")
	}
	if testutil.NewTestRepoAt(abcdef, "/a/b/c/d/e/f") == nil {
		panic("NewTestRepoAt() returned nil")
	}

	other := testlib.TempDirAbs(t)

	otherabc := filepath.Join(other, "a", "b", "c")
	if testutil.NewTestRepoAt(otherabc, "/a/b/c") == nil {
		panic("NewTestRepoAt() returned nil")
	}

	type args struct {
		clonePath string
	}
	tests := []struct {
		name         string
		patternValue string
		args         args
		want         float64
	}{
		// matching the empty pattern
		{"EmptyPattern", "", args{abc}, 1},

		// matching one-component parts of a/b/c
		{"oneComponentStart", "a", args{abc}, 0.25},
		{"oneComponentStart outside root", "a", args{otherabc}, 0.25},
		{"oneComponentMiddle", "b", args{abc}, 0.5},
		{"oneComponentEnd", "c", args{abc}, 1},
		{"oneComponentNot", "d", args{abc}, -1},

		// matching constant sub-paths
		{"twoComponentsConst", "b/c", args{abcdef}, 0.125},
		{"noTwoComponentsConst", "f/g", args{abcdef}, -1},

		// variable sub-paths
		{"variableSubPathPositive", "b/*/d", args{abcdef}, 0.25},
		{"variableSubPathNegative", "b/*/c", args{abcdef}, -1},
	}
	for _, tt := range tests {
		var pat PatternFilter
		t.Run(tt.name, func(t *testing.T) {
			pat.Set(tt.patternValue)
			if got := pat.Score(
				Env{
					Root: root,
					Git:  git.NewGitFromPlumbing(nil, ""),
				},
				path.ToOSPath(tt.args.clonePath),
			); got != tt.want {
				t.Errorf("PatternFilter().Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatternFilter_MatchesURL(t *testing.T) {
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
		var pat PatternFilter
		t.Run(tt.name, func(t *testing.T) {
			pat.Set(tt.args.pattern)
			if got := pat.MatchesURL(ParseURL(tt.args.s)); got != tt.want {
				t.Errorf("PatternFilter().MatchesString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisjunctionFilter_Score(t *testing.T) {
	type fields struct {
		Clauses []Filter
	}
	type args struct {
		root      string
		clonePath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			"zero filters never match",
			fields{
				Clauses: nil,
			},
			args{
				"/root/",
				"/root/whatever/",
			},
			-1,
		},

		{
			"two pathfilters match first path",
			fields{
				Clauses: []Filter{
					PathFilter{[]string{path.ToOSPath("/root/matcha")}},
					PathFilter{[]string{path.ToOSPath("/root/matchb")}},
				},
			},
			args{
				"/root/",
				"/root/matcha",
			},
			1,
		},

		{
			"two pathfilters match second path",
			fields{
				Clauses: []Filter{
					PathFilter{[]string{path.ToOSPath("/root/matcha")}},
					PathFilter{[]string{path.ToOSPath("/root/matchb")}},
				},
			},
			args{
				"/root/",
				"/root/matchb",
			},
			1,
		},

		{
			"two pathfilters do not match third path",
			fields{
				Clauses: []Filter{
					PathFilter{[]string{path.ToOSPath("/root/matcha")}},
					PathFilter{[]string{path.ToOSPath("/root/matchb")}},
				},
			},
			args{
				"/root/",
				"/root/matchc",
			},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			or := DisjunctionFilter{
				Clauses: tt.fields.Clauses,
			}
			if got := or.Score(
				Env{Root: path.ToOSPath(tt.args.root)},
				path.ToOSPath(tt.args.clonePath),
			); got != tt.want {
				t.Errorf("DisjunctionFilter.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisjunctionFilter_Candidates(t *testing.T) {
	type fields struct {
		Clauses []Filter
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"zero filters don't have candidates",
			fields{
				Clauses: nil,
			},
			[]string{},
		},

		{
			"two candidates get returned",
			fields{
				Clauses: []Filter{
					PathFilter{[]string{path.ToOSPath("/root/matcha")}},
					PathFilter{[]string{path.ToOSPath("/root/matchb")}},
				},
			},
			path.ToOSPaths([]string{
				"/root/matcha",
				"/root/matchb",
			}),
		},

		{
			"duplicate candidates get returned only once",
			fields{
				Clauses: []Filter{
					PathFilter{[]string{path.ToOSPath("/root/matcha")}},
					PathFilter{[]string{path.ToOSPath("/root/matchb")}},
					PathFilter{[]string{path.ToOSPath("/root/matchb")}},
				},
			},
			path.ToOSPaths([]string{
				"/root/matcha",
				"/root/matchb",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			or := DisjunctionFilter{
				Clauses: tt.fields.Clauses,
			}
			if got := or.Candidates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DisjunctionFilter.Candidates() = %v, want %v", got, tt.want)
			}
		})
	}
}
