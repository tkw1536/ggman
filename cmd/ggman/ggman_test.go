package main

import (
	"testing"

	"github.com/tkw1536/ggman/internal/mockenv"
	"github.com/tkw1536/goprogram/lib/doccheck"
)

// This test runs every command once with the --help flag
//
// This tests that the Description() functions do not fail.
// This also checks that all the docstrings are valid if the doccheck flag is specified.
func Test_main_docs(t *testing.T) {
	mock := mockenv.NewMockEnv(t)

	for _, name := range ggmanExe.Commands() {
		t.Run(name, func(t *testing.T) {
			defer func() {
				err := recover()
				if ve, isVe := err.(*doccheck.CheckError); isVe {
					t.Fatalf("Doccheck failure: %s", ve)
				}
			}()
			cmd, ok := ggmanExe.Command(name)
			if !ok {
				t.Fail()
			}
			code, _, _ := mock.Run(cmd, "", "", name, "--help")
			if code != 0 {
				t.Fail()
			}
		})
	}
}
