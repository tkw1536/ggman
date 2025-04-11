//spellchecker:words path
package path_test

//spellchecker:words path filepath testing github ggman internal testutil pkglib testlib
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/ggman/internal/path"
	"github.com/tkw1536/ggman/internal/testutil"
	"github.com/tkw1536/pkglib/testlib"
)

func TestJoinNormalized(t *testing.T) {
	t.Parallel()

	// create subdirectory for testing
	root := testlib.TempDirAbs(t)
	exactD := filepath.Join(root, "exact")
	if err := os.Mkdir(exactD, os.ModePerm|os.ModeDir); err != nil {
		panic(err)
	}

	type args struct {
		n    path.Normalization
		base string
		elem []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"join new directory (NoNorm)",
			args{
				n:    path.NoNorm,
				base: root,
				elem: []string{"new"},
			},
			filepath.Join(root, "new"),
			false,
		},
		{
			"join new directory (FoldNorm)",
			args{
				n:    path.FoldNorm,
				base: root,
				elem: []string{"new"},
			},
			filepath.Join(root, "new"),
			false,
		},
		{
			"join new directory (FoldPreferExactNorm)",
			args{
				n:    path.FoldPreferExactNorm,
				base: root,
				elem: []string{"new"},
			},
			filepath.Join(root, "new"),
			false,
		},

		// join exact directory
		{
			"join exact directory (NoNorm)",
			args{
				n:    path.NoNorm,
				base: root,
				elem: []string{"exact"},
			},
			filepath.Join(root, "exact"),
			false,
		},
		{
			"join exact directory (FoldNorm)",
			args{
				n:    path.FoldNorm,
				base: root,
				elem: []string{"exact"},
			},
			filepath.Join(root, "exact"),
			false,
		},
		{
			"join exact directory (FoldPreferExactNorm)",
			args{
				n:    path.FoldPreferExactNorm,
				base: root,
				elem: []string{"exact"},
			},
			filepath.Join(root, "exact"),
			false,
		},

		// join in-exact match
		{
			"join inexact directory (NoNorm)",
			args{
				n:    path.NoNorm,
				base: root,
				elem: []string{"eXact"},
			},
			filepath.Join(root, "eXact"),
			false,
		},
		{
			"join inexact directory (FoldNorm)",
			args{
				n:    path.FoldNorm,
				base: root,
				elem: []string{"eXact"},
			},
			filepath.Join(root, "exact"),
			false,
		},
		{
			"join inexact directory (FoldPreferExactNorm)",
			args{
				n:    path.FoldPreferExactNorm,
				base: root,
				elem: []string{"eXact"},
			},
			filepath.Join(root, "exact"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := path.JoinNormalized(tt.args.n, tt.args.base, tt.args.elem...)
			if (err != nil) != tt.wantErr {
				t.Errorf("JoinNormalized() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JoinNormalized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoinNormalized_both(t *testing.T) {
	t.Parallel()

	if !testutil.CaseSensitive(t) {
		t.Skipf("Filesystem is case-insensitive")
	}

	// create subdirectory for testing
	root := testlib.TempDirAbs(t)

	lcBothD := filepath.Join(root, "both")
	if err := os.Mkdir(lcBothD, os.ModePerm|os.ModeDir); err != nil {
		panic(err)
	}

	ucBothD := filepath.Join(root, "BOTH")
	if err := os.Mkdir(ucBothD, os.ModePerm|os.ModeDir); err != nil {
		panic(err)
	}

	type args struct {
		n    path.Normalization
		base string
		elem []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// 'BOTH'
		{
			"join first existing directory (NoNorm)",
			args{
				n:    path.NoNorm,
				base: root,
				elem: []string{"BOTH"},
			},
			filepath.Join(root, "BOTH"),
			false,
		},
		{
			"join first existing directory (FoldNorm)",
			args{
				n:    path.FoldNorm,
				base: root,
				elem: []string{"BOTH"},
			},
			filepath.Join(root, "BOTH"),
			false,
		},
		{
			"join first existing directory (FoldPreferExactNorm)",
			args{
				n:    path.FoldPreferExactNorm,
				base: root,
				elem: []string{"BOTH"},
			},
			filepath.Join(root, "BOTH"),
			false,
		},

		// 'both'
		{
			"join second existing directory (NoNorm)",
			args{
				n:    path.NoNorm,
				base: root,
				elem: []string{"both"},
			},
			filepath.Join(root, "both"),
			false,
		},
		{
			"join second existing directory (FoldNorm)",
			args{
				n:    path.FoldNorm,
				base: root,
				elem: []string{"both"},
			},
			filepath.Join(root, "BOTH"),
			false,
		},
		{
			"join second existing directory (FoldPreferExactNorm)",
			args{
				n:    path.FoldPreferExactNorm,
				base: root,
				elem: []string{"both"},
			},
			filepath.Join(root, "both"),
			false,
		},

		// 'BoTh'
		{
			"join neither existing directory (NoNorm)",
			args{
				n:    path.NoNorm,
				base: root,
				elem: []string{"BoTh"},
			},
			filepath.Join(root, "BoTh"),
			false,
		},
		{
			"join neither existing directory (FoldNorm)",
			args{
				n:    path.FoldNorm,
				base: root,
				elem: []string{"BoTh"},
			},
			filepath.Join(root, "BOTH"),
			false,
		},
		{
			"join neither existing directory (FoldPreferExactNorm)",
			args{
				n:    path.FoldPreferExactNorm,
				base: root,
				elem: []string{"BoTh"},
			},
			filepath.Join(root, "BOTH"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := path.JoinNormalized(tt.args.n, tt.args.base, tt.args.elem...)
			if (err != nil) != tt.wantErr {
				t.Errorf("JoinNormalized() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JoinNormalized() = %v, want %v", got, tt.want)
			}
		})
	}
}
