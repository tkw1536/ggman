package stream

// FromNil creates a new IOStream that silences all output and provides no input.
func FromNil() IOStream {
	return NewIOStream(nil, nil, nil, 0)
}
