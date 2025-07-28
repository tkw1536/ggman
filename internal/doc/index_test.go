package doc_test

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/doc"
)

func TestIndexFilenames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupCmd  func() *cobra.Command
		extension string
		want      map[string]string
	}{
		{
			name: "single command with md extension",
			setupCmd: func() *cobra.Command {
				return &cobra.Command{Use: "root"}
			},
			extension: "md",
			want: map[string]string{
				"root.md": "root",
			},
		},
		{
			name: "command with children and txt extension",
			setupCmd: func() *cobra.Command {
				root := &cobra.Command{Use: "root"}
				child1 := &cobra.Command{Use: "child1"}
				child2 := &cobra.Command{Use: "child2"}
				root.AddCommand(child1, child2)
				return root
			},
			extension: "txt",
			want: map[string]string{
				"root.txt":        "root",
				"root_child1.txt": "root child1",
				"root_child2.txt": "root child2",
			},
		},
		{
			name: "nested command tree with html extension",
			setupCmd: func() *cobra.Command {
				root := &cobra.Command{Use: "root"}
				child1 := &cobra.Command{Use: "child1"}
				child2 := &cobra.Command{Use: "child2"}
				grandchild1 := &cobra.Command{Use: "grandchild1"}
				grandchild2 := &cobra.Command{Use: "grandchild2"}

				child1.AddCommand(grandchild1)
				child2.AddCommand(grandchild2)
				root.AddCommand(child1, child2)
				return root
			},
			extension: "html",
			want: map[string]string{
				"root.html":                    "root",
				"root_child1.html":             "root child1",
				"root_child1_grandchild1.html": "root child1 grandchild1",
				"root_child2.html":             "root child2",
				"root_child2_grandchild2.html": "root child2 grandchild2",
			},
		},
		{
			name: "empty extension",
			setupCmd: func() *cobra.Command {
				root := &cobra.Command{Use: "root"}
				child := &cobra.Command{Use: "child"}
				root.AddCommand(child)
				return root
			},
			extension: "",
			want: map[string]string{
				"root.":       "root",
				"root_child.": "root child",
			},
		},
		{
			name: "complex extension",
			setupCmd: func() *cobra.Command {
				root := &cobra.Command{Use: "root"}
				child := &cobra.Command{Use: "child"}
				root.AddCommand(child)
				return root
			},
			extension: "markdown",
			want: map[string]string{
				"root.markdown":       "root",
				"root_child.markdown": "root child",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cmd := tt.setupCmd()
			got := doc.IndexFilenames(cmd, tt.extension)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexFilenames() = %v, want %v", got, tt.want)
			}
		})
	}
}
