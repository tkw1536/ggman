package program

import (
	"reflect"
	"testing"
)

func TestProgram_Alias(t *testing.T) {
	var program tProgram

	program.Register(fakeCommand("a"))
	program.RegisterAlias(Alias{Name: "a", Command: "b", Args: []string{"c"}})
	program.RegisterAlias(Alias{Name: "b", Command: "d", Args: []string{"e"}})

	got := program.Aliases()
	want := []string{"a", "b"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Program.Aliases() = %v, want = %v", got, want)
	}
}
