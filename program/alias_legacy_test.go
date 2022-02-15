package program_test

import (
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/program"
)

func TestProgram_Alias(t *testing.T) {
	var p tProgram

	p.Register(fakeCommand("a"))
	p.RegisterAlias(program.Alias{Name: "a", Command: "b", Args: []string{"c"}})
	p.RegisterAlias(program.Alias{Name: "b", Command: "d", Args: []string{"e"}})

	got := p.Aliases()
	want := []string{"a", "b"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Program.Aliases() = %v, want = %v", got, want)
	}
}
