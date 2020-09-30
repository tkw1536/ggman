package util

import "io"

// NullStream is an interface that implements io.Reader and io.Writer as no-ops.
type NullStream struct{}

// Read is a no-op that always returns 0, io.EOF.
func (NullStream) Read(bytes []byte) (int, error) {
	return 0, io.EOF
}

// Write is a no-op that always returns len(bytes), nil.
func (NullStream) Write(bytes []byte) (int, error) {
	return len(bytes), nil
}

// Close is a no-op that always returns nil.
func (NullStream) Close() error {
	return nil
}
