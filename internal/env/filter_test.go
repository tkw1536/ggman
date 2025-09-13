package env_test

//spellchecker:words context path filepath reflect testing ggman internal testutil pkglib testlib
import (
	"context"
	"path/filepath"
	"reflect"
	"testing"

	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/ggman/internal/git"
	"go.tkw01536.de/ggman/internal/testutil"
	"go.tkw01536.de/pkglib/testlib"
)

//spellchecker:words matcha matchb matchc otherabc

func setupFilterTest(t *testing.T) (root, exampleClonePath, otherClonePath string) {
	t.Helper()

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

type testFilter struct{}

func (testFilter) Score(ctx context.Context, env *env.Env, clonePath string) float64 {
	panic("never reached")
}

type testFilterWithCandidates struct {
	testFilter
}

func (testFilterWithCandidates) Candidates() []string { return []string{"a", "b", "c"} }

func TestCandidates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		filter env.Filter
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
			t.Parallel()

			if got := env.Candidates(tt.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Candidates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathFilter_Score(t *testing.T) {
	t.Parallel()

	root, exampleClonePath, otherClonePath := setupFilterTest(t)

	type fields struct {
		Paths []string
	}
	type args struct {
		env       *env.Env
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
				&env.Env{Root: root},
				root,
			},
			env.FilterDoesNotMatch,
		},
		{
			"non-listed path doesn't match",
			fields{
				Paths: []string{exampleClonePath, otherClonePath},
			},
			args{
				&env.Env{Root: root},
				"/outside/",
			},
			env.FilterDoesNotMatch,
		},
		{
			"listed path matches (1)",
			fields{
				Paths: []string{exampleClonePath, otherClonePath},
			},
			args{
				&env.Env{Root: root},
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
				&env.Env{Root: root},
				otherClonePath,
			},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pf := env.PathFilter{
				Paths: tt.fields.Paths,
			}
			if got := pf.Score(t.Context(), tt.args.env, tt.args.clonePath); got != tt.want {
				t.Errorf("PathFilter.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathFilter_Candidates(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			pf := env.PathFilter{
				Paths: tt.fields.Paths,
			}
			if got := pf.Candidates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PathFilter.Candidates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatternFilter_Score(t *testing.T) {
	t.Parallel()

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
		{"oneComponentStart", "a", args{abc}, 0.5},
		{"oneComponentStart outside root", "a", args{otherabc}, 0.5},
		{"oneComponentMiddle", "b", args{abc}, 0.5},
		{"oneComponentEnd", "c", args{abc}, 1},
		{"oneComponentNot", "d", args{abc}, env.FilterDoesNotMatch},

		// matching constant sub-paths
		{"twoComponentsConst", "b/c", args{abcdef}, 0.25},
		{"noTwoComponentsConst", "f/g", args{abcdef}, env.FilterDoesNotMatch},

		// variable sub-paths
		{"variableSubPathPositive", "b/*/d", args{abcdef}, 0.25},
		{"variableSubPathNegative", "b/*/c", args{abcdef}, env.FilterDoesNotMatch},
	}
	for _, tt := range tests {
		var pat env.PatternFilter
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pat.Set(tt.patternValue)
			if got := pat.Score(
				t.Context(),
				&env.Env{
					Root: root,
					Git:  git.NewGitFromPlumbing(nil, ""),
				},
				testutil.ToOSPath(tt.args.clonePath),
			); got != tt.want {
				t.Errorf("PatternFilter().Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatternFilter_MatchesURL(t *testing.T) {
	t.Parallel()

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
		var pat env.PatternFilter
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pat.Set(tt.args.pattern)
			if got := pat.MatchesURL(env.ParseURL(tt.args.s)); got != tt.want {
				t.Errorf("PatternFilter().MatchesString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisjunctionFilter_Score(t *testing.T) {
	t.Parallel()

	type fields struct {
		Clauses []env.Filter
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
			env.FilterDoesNotMatch,
		},

		{
			"two PathFilters match first path",
			fields{
				Clauses: []env.Filter{
					env.PathFilter{[]string{testutil.ToOSPath("/root/matcha")}},
					env.PathFilter{[]string{testutil.ToOSPath("/root/matchb")}},
				},
			},
			args{
				"/root/",
				"/root/matcha",
			},
			1,
		},

		{
			"two PathFilters match second path",
			fields{
				Clauses: []env.Filter{
					env.PathFilter{[]string{testutil.ToOSPath("/root/matcha")}},
					env.PathFilter{[]string{testutil.ToOSPath("/root/matchb")}},
				},
			},
			args{
				"/root/",
				"/root/matchb",
			},
			1,
		},

		{
			"two PathFilters do not match third path",
			fields{
				Clauses: []env.Filter{
					env.PathFilter{[]string{testutil.ToOSPath("/root/matcha")}},
					env.PathFilter{[]string{testutil.ToOSPath("/root/matchb")}},
				},
			},
			args{
				"/root/",
				"/root/matchc",
			},
			env.FilterDoesNotMatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			or := env.DisjunctionFilter{
				Clauses: tt.fields.Clauses,
			}
			if got := or.Score(
				t.Context(),
				&env.Env{Root: testutil.ToOSPath(tt.args.root)},
				testutil.ToOSPath(tt.args.clonePath),
			); got != tt.want {
				t.Errorf("DisjunctionFilter.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisjunctionFilter_Candidates(t *testing.T) {
	t.Parallel()

	type fields struct {
		Clauses []env.Filter
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
				Clauses: []env.Filter{
					env.PathFilter{[]string{testutil.ToOSPath("/root/matcha")}},
					env.PathFilter{[]string{testutil.ToOSPath("/root/matchb")}},
				},
			},
			testutil.ToOSPaths([]string{
				"/root/matcha",
				"/root/matchb",
			}),
		},

		{
			"duplicate candidates get returned only once",
			fields{
				Clauses: []env.Filter{
					env.PathFilter{[]string{testutil.ToOSPath("/root/matcha")}},
					env.PathFilter{[]string{testutil.ToOSPath("/root/matchb")}},
					env.PathFilter{[]string{testutil.ToOSPath("/root/matchb")}},
				},
			},
			testutil.ToOSPaths([]string{
				"/root/matcha",
				"/root/matchb",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			or := env.DisjunctionFilter{
				Clauses: tt.fields.Clauses,
			}
			if got := or.Candidates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DisjunctionFilter.Candidates() = %v, want %v", got, tt.want)
			}
		})
	}
}
