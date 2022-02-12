package stream

import (
	"io"
	"testing"
)

func TestNullStream_Read(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantEOF bool
	}{
		{"read 0 bytes into nil slice", args{nil}, 0, true},
		{"read 0 bytes into non-nil slice", args{make([]byte, 10)}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Null.Read(tt.args.bytes)
			if (err == io.EOF) != tt.wantEOF {
				t.Errorf("NullStream.Read() error = %v, wantEOF %v", err, tt.wantEOF)
				return
			}
			if got != tt.want {
				t.Errorf("NullStream.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullStream_Write(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"write 0 bytes into nil slice", args{nil}, 0, false},
		{"write 10 bytes into non-nil slice", args{make([]byte, 10)}, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Null.Write(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("NullStream.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NullStream.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullStream_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"return nil", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Null.Close(); (err != nil) != tt.wantErr {
				t.Errorf("NullStream.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
