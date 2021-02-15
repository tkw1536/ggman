package testutil

import (
	"os"
	"testing"
)

func TestMockVariables(t *testing.T) {
	os.Setenv("A", "oldA")
	os.Setenv("B", "")
	os.Setenv("C", "oldC")

	undo := MockVariables(map[string]string{
		"A": "newA",
		"B": "newB",
		"C": "newC",
	})

	if os.Getenv("A") != "newA" {
		t.Errorf("MockVariables() did not set A")
	}
	if os.Getenv("B") != "newB" {
		t.Errorf("MockVariables() did not set B")
	}
	if os.Getenv("C") != "newC" {
		t.Errorf("MockVariables() did not set C")
	}

	undo()

	if os.Getenv("A") != "oldA" {
		t.Errorf("MockVariables() did not revert A")
	}
	if os.Getenv("B") != "" {
		t.Errorf("MockVariables() did not revert B")
	}
	if os.Getenv("C") != "oldC" {
		t.Errorf("MockVariables() did not revert C")
	}
}
