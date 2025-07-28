package doc_test

//spellchecker:words testing github cobra ggman internal
import (
	"testing"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/doc"
)

func TestAllCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupCmd func() *cobra.Command
		want     []string
	}{
		{
			name: "single command",
			setupCmd: func() *cobra.Command {
				return &cobra.Command{Use: "root"}
			},
			want: []string{"root"},
		},
		{
			name: "command with children",
			setupCmd: func() *cobra.Command {
				root := &cobra.Command{Use: "root"}
				child1 := &cobra.Command{Use: "child1"}
				child2 := &cobra.Command{Use: "child2"}
				root.AddCommand(child1, child2)
				return root
			},
			want: []string{"root", "child1", "child2"},
		},
		{
			name: "nested command tree",
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
			want: []string{"root", "child1", "grandchild1", "child2", "grandchild2"},
		},
		{
			name: "deep nested tree",
			setupCmd: func() *cobra.Command {
				root := &cobra.Command{Use: "root"}
				level1 := &cobra.Command{Use: "level1"}
				level2 := &cobra.Command{Use: "level2"}
				level3 := &cobra.Command{Use: "level3"}
				level4 := &cobra.Command{Use: "level4"}

				level3.AddCommand(level4)
				level2.AddCommand(level3)
				level1.AddCommand(level2)
				root.AddCommand(level1)
				return root
			},
			want: []string{"root", "level1", "level2", "level3", "level4"},
		},
		{
			name: "multiple children at same level",
			setupCmd: func() *cobra.Command {
				root := &cobra.Command{Use: "root"}
				child1 := &cobra.Command{Use: "child1"}
				child2 := &cobra.Command{Use: "child2"}
				child3 := &cobra.Command{Use: "child3"}

				root.AddCommand(child1, child2, child3)
				return root
			},
			want: []string{"root", "child1", "child2", "child3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cmd := tt.setupCmd()
			var got []string

			for c := range doc.AllCommands(cmd) {
				got = append(got, c.Use)
			}

			if len(got) != len(tt.want) {
				t.Errorf("AllCommands() returned %d commands, want %d", len(got), len(tt.want))
				return
			}

			for i, want := range tt.want {
				if got[i] != want {
					t.Errorf("AllCommands()[%d] = %s, want %s", i, got[i], want)
				}
			}
		})
	}
}

func TestCountCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupCmd func() *cobra.Command
		want     int
	}{
		{
			name: "single command",
			setupCmd: func() *cobra.Command {
				return &cobra.Command{Use: "root"}
			},
			want: 1,
		},
		{
			name: "command with two children",
			setupCmd: func() *cobra.Command {
				root := &cobra.Command{Use: "root"}
				child1 := &cobra.Command{Use: "child1"}
				child2 := &cobra.Command{Use: "child2"}
				root.AddCommand(child1, child2)
				return root
			},
			want: 3,
		},
		{
			name: "complex tree",
			setupCmd: func() *cobra.Command {
				root := &cobra.Command{Use: "root"}
				child1 := &cobra.Command{Use: "child1"}
				child2 := &cobra.Command{Use: "child2"}
				grandchild1 := &cobra.Command{Use: "grandchild1"}
				grandchild2 := &cobra.Command{Use: "grandchild2"}
				grandchild3 := &cobra.Command{Use: "grandchild3"}

				child1.AddCommand(grandchild1, grandchild2)
				child2.AddCommand(grandchild3)
				root.AddCommand(child1, child2)
				return root
			},
			want: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cmd := tt.setupCmd()
			got := doc.CountCommands(cmd)

			if got != tt.want {
				t.Errorf("CountCommands() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAllCommandsEarlyReturn(t *testing.T) {
	t.Parallel()

	// Test that the iterator stops when yield returns false
	root := &cobra.Command{Use: "root"}
	child1 := &cobra.Command{Use: "child1"}
	child2 := &cobra.Command{Use: "child2"}
	child3 := &cobra.Command{Use: "child3"}
	root.AddCommand(child1, child2, child3)

	visited := make([]string, 0, 2)
	count := 0

	for c := range doc.AllCommands(root) {
		visited = append(visited, c.Use)
		count++
		if count == 2 {
			break // This should stop the iteration
		}
	}

	if len(visited) != 2 {
		t.Errorf("Expected to visit 2 commands, but visited %d", len(visited))
	}

	expected := []string{"root", "child1"}
	for i, want := range expected {
		if visited[i] != want {
			t.Errorf("visited[%d] = %s, want %s", i, visited[i], want)
		}
	}
}
