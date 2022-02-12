// Package stream provides NullStream
package stream

import "io"

// Null is an io.ReadWriteCloser.
//
// Reads from it return 0 bytes and io.EOF.
// Writes and Closes succeed without doing anything.
//
// See also io.Discard.
var Null io.ReadWriteCloser = nullStream{}

type nullStream struct{}

func (nullStream) Read(bytes []byte) (int, error)  { return 0, io.EOF }
func (nullStream) Write(bytes []byte) (int, error) { return len(bytes), nil }
func (nullStream) Close() error                    { return nil }
