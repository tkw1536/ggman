package exit

import (
	"errors"
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/internal/testutil"
)

func TestAsError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want Error
	}{
		{"nil error returns zero value", args{nil}, Error{}},
		{"Error object returns itself", args{Error{ExitCode: ExitGeneric, Message: "stuff"}}, Error{ExitCode: ExitGeneric, Message: "stuff"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsErrorPanic(t *testing.T) {
	_, gotPanic := testutil.DoesPanic(func() { AsError(errors.New("not an error")) })
	wantPanic := interface{}("AsError: err must be nil or Error")
	if wantPanic != gotPanic {
		t.Errorf("AsError: got panic = %v, want = %v", gotPanic, wantPanic)
	}
}

func TestError_WithMessage(t *testing.T) {
	type fields struct {
		ExitCode ExitCode
		Message  string
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Error
	}{
		{"replaces empty message", fields{}, args{message: "Hello world"}, Error{Message: "Hello world"}},
		{"replaces non-empty message", fields{Message: "not empty"}, args{message: "Hello world"}, Error{Message: "Hello world"}},

		{"keeps exit code 1", fields{ExitCode: 1}, args{message: "Hello world"}, Error{ExitCode: 1, Message: "Hello world"}},
		{"keeps exit code 2", fields{ExitCode: 2}, args{message: "Hello world"}, Error{ExitCode: 2, Message: "Hello world"}},
		{"keeps exit code 3", fields{ExitCode: 3}, args{message: "Hello world"}, Error{ExitCode: 3, Message: "Hello world"}},
		{"keeps exit code 4", fields{ExitCode: 4}, args{message: "Hello world"}, Error{ExitCode: 4, Message: "Hello world"}},
		{"keeps exit code 5", fields{ExitCode: 5}, args{message: "Hello world"}, Error{ExitCode: 5, Message: "Hello world"}},

		{"does not substitute strings in old message", fields{Message: "old %s"}, args{message: "Hello world"}, Error{Message: "Hello world"}},
		{"does not substitute strings in new message", fields{Message: "old message"}, args{message: "Hello world %s"}, Error{Message: "Hello world %s"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Error{
				ExitCode: tt.fields.ExitCode,
				Message:  tt.fields.Message,
			}
			if got := err.WithMessage(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.WithMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WithMessageF(t *testing.T) {
	type fields struct {
		ExitCode ExitCode
		Message  string
	}
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Error
	}{
		{"keeps message without format", fields{Message: "Hello world"}, args{}, Error{Message: "Hello world"}},
		{"replaces message", fields{Message: "Hello %s"}, args{[]interface{}{"world"}}, Error{Message: "Hello world"}},

		{"keeps exit code 1", fields{ExitCode: 1, Message: "%s"}, args{[]interface{}{"Hello world"}}, Error{ExitCode: 1, Message: "Hello world"}},
		{"keeps exit code 2", fields{ExitCode: 2, Message: "%s"}, args{[]interface{}{"Hello world"}}, Error{ExitCode: 2, Message: "Hello world"}},
		{"keeps exit code 3", fields{ExitCode: 3, Message: "%s"}, args{[]interface{}{"Hello world"}}, Error{ExitCode: 3, Message: "Hello world"}},
		{"keeps exit code 4", fields{ExitCode: 4, Message: "%s"}, args{[]interface{}{"Hello world"}}, Error{ExitCode: 4, Message: "Hello world"}},
		{"keeps exit code 5", fields{ExitCode: 5, Message: "%s"}, args{[]interface{}{"Hello world"}}, Error{ExitCode: 5, Message: "Hello world"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Error{
				ExitCode: tt.fields.ExitCode,
				Message:  tt.fields.Message,
			}
			if got := err.WithMessageF(tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.WithMessageF() = %v, want %v", got, tt.want)
			}
		})
	}
}
