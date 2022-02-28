package goprogram

import (
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/goprogram/stream"
)

// Register an alias for a program.
// See the test suite for instaniated types.
func ExampleAlias() {
	// create a new program that only has an echo command
	// this code is reused across the test suite, hence not shown here.
	p := makeProgram()
	p.Register(makeEchoCommand("echo"))

	// register an alias "hello" which expands into hello world
	p.RegisterAlias(Alias{Name: "hello", Command: "echo", Args: []string{"hello world"}})

	p.Main(stream.FromEnv(), "", []string{"hello"})
	p.Main(stream.FromEnv(), "", []string{"hello", "again"})

	// Output: [hello world]
	// [hello world again]
}

func TestProgram_Alias(t *testing.T) {
	var p iProgram

	p.Register(tCommand{desc: iDescription{
		Command: "a",
	}})
	p.RegisterAlias(Alias{Name: "a", Command: "b", Args: []string{"c"}})
	p.RegisterAlias(Alias{Name: "b", Command: "d", Args: []string{"e"}})

	got := p.Aliases()
	want := []string{"a", "b"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Program.Aliases() = %v, want = %v", got, want)
	}
}
