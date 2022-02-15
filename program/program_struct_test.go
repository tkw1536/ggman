package program

import (
	"fmt"

	"github.com/tkw1536/ggman/program/meta"
)

// This file contains dummy implementations of everything required to assemble a program.
// It is reused across the test suite, however there is no versioning guarantee.
// It may change in a future revision of the test suite.

// Runtime of each command is a single string value.
// Parameters to initialize each command is also a string value.
type ttRuntime string
type ttParameters string

// tFlags holds a set of dummy global flags.
type ttFlags struct {
}

// tRequirements is the implementation of the AllowsFlag function
type ttRequirements func(flag meta.Flag) bool

func (t ttRequirements) AllowsFlag(flag meta.Flag) bool { return t(flag) }
func (t ttRequirements) Validate(args Arguments[ttFlags]) error {
	return ValidateAllowedFlags[ttFlags](t, args)
}

// instiantiated types for the test suite
type iProgram = Program[ttRuntime, ttParameters, ttFlags, ttRequirements]
type iCommand = Command[ttRuntime, ttParameters, ttFlags, ttRequirements]
type iContext = Context[ttRuntime, ttParameters, ttFlags, ttRequirements]
type iArguments = Arguments[ttFlags]
type iDescription = Description[ttFlags, ttRequirements]

// ttCommand represents a sample test suite command.
// It runs the associated private functions, or prints an info message to stdout.
type ttCommand[F any] struct {
	flags F `group:"F"` // flags holds command flags

	beforeRegister func() error
	desc           iDescription
	afterParse     func() error
	run            func(context iContext) error
}

func (t ttCommand[F]) BeforeRegister(program *iProgram) {
	if t.beforeRegister == nil {
		fmt.Println("BeforeRegister()")
		return
	}
	t.beforeRegister()
}
func (t ttCommand[F]) Description() iDescription {
	return t.desc
}
func (t ttCommand[F]) AfterParse() error {
	if t.afterParse == nil {
		fmt.Println("AfterParse()")
		return nil
	}
	return t.afterParse()
}
func (t ttCommand[F]) Run(ctx iContext) error {
	if t.run == nil {
		fmt.Println("Run()")
		return nil
	}
	return t.run(ctx)
}
