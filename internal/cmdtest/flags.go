// Package cmdtest provides common command testing utilities.
//
//spellchecker:words cmdtest
package cmdtest

//spellchecker:words reflect slices ggman internal testutil goprogram meta parser pkglib reflectx
import (
	"reflect"
	"slices"

	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/internal/testutil"
	"go.tkw01536.de/goprogram"
	"go.tkw01536.de/goprogram/meta"
	"go.tkw01536.de/goprogram/parser"
	"go.tkw01536.de/pkglib/reflectx"
)

// AssertNoFlagOverlap asserts that there is no overlap between the flags for command and ggman global flags.
func AssertNoFlagOverlap(t testutil.TestingT, command ggman.Command) {
	t.Helper()

	assertFlagOverlap(t, command, []string{})
}

func assertFlagOverlap(t testutil.TestingT, command ggman.Command, want []string) {
	t.Helper()

	cCommand, _ := reflectx.CopyInterface(command)
	got := flagOverlap(cCommand)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got FlagOverlap = %v, but wanted %v", got, want)
	}
}

func flagOverlap[E any, P any, F any, R goprogram.Requirement[F]](command goprogram.Command[E, P, F, R]) []string {
	globals := flagNames(parser.AllFlags[F]())
	locals := flagNames(parser.AllFlagsOf(command))

	globalSet := make(map[string]struct{}, len(globals))
	for _, f := range globals {
		globalSet[f] = struct{}{}
	}

	// determine the overlap
	overlap := make([]string, 0, min(len(locals), len(globalSet)))
	for _, f := range locals {
		if _, ok := globalSet[f]; ok {
			overlap = append(overlap, f)
		}
	}
	slices.Sort(overlap)
	return overlap
}

func flagNames(flags []meta.Flag) []string {
	// cap: typically we have one short and one long flag
	names := make([]string, 0, 2*len(flags))
	for _, f := range flags {
		names = append(names, f.Long...)
		names = append(names, f.Short...)
	}
	return names
}
