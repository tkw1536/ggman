package goprogram

import (
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/goprogram/stream"
)

// Register a command for a program.
// See the test suite for instaniated types.
func ExampleCommand() {
	// create a new program that only has an echo command
	// this code is reused across the test suite, hence not shown here.
	p := makeProgram()
	p.Register(makeEchoCommand("echo"))

	// Execute the command with some arguments
	p.Main(stream.FromEnv(), "", []string{"echo", "hello", "world"})

	// Output: [hello world]
}

func TestProgram_Commands(t *testing.T) {
	p := makeProgram()
	p.Register(makeEchoCommand("a"))
	p.Register(makeEchoCommand("c"))
	p.Register(makeEchoCommand("b"))

	got := p.Commands()
	want := []string{"a", "b", "c"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Program.Commands() = %v, want = %v", got, want)
	}
}

func TestProgram_FmtCommands(t *testing.T) {
	p := makeProgram()
	p.Register(makeEchoCommand("a"))
	p.Register(makeEchoCommand("c"))
	p.Register(makeEchoCommand("b"))

	got := p.FmtCommands()
	want := `"a", "b", "c"`

	if got != want {
		t.Errorf("Program.FmtCommands() = %v, want = %v", got, want)
	}
}
