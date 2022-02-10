package exit

import "fmt"

// Error represents any error state by a program.
// It implements the builtin error interface.
//
// The zero value represents that no error occured and is ready to use.
type Error struct {
	// Exit code of the program (if applicable)
	ExitCode
	// Message for this error
	Message string
}

// AsError asserts that err is either nil or of type Error and returns it.
// When err is nil, the zero value of type Error is returned.
//
// If err is not nil and not of type Error, calls panic().
func AsError(err error) Error {
	switch e := err.(type) {
	case nil:
		return Error{}
	case Error:
		return e
	}
	panic("AsError: err must be nil or Error")
}

// WithMessage returns a copy of this error with the same Code but different Message.
// The new message is the message passed as an argument.
func (err Error) WithMessage(message string) Error {
	return Error{
		ExitCode: err.ExitCode,
		Message:  message,
	}
}

// WithMessageF returns a copy of this error with the same Code but different Message.
// The new message is the current message, formatted by SPrintf and the arguments.
func (err Error) WithMessageF(args ...interface{}) Error {
	return err.WithMessage(fmt.Sprintf(err.Message, args...))
}

func init() {
	var _ error = (*Error)(nil)
}

func (err Error) Error() string {
	return err.Message
}
