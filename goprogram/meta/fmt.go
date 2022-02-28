// Package meta contains facilities to provide meta-information about programs and commands.
package meta

import (
	"strings"
	"sync"
)

// builderPool used by various formatters in this package
var builderPool = &sync.Pool{
	New: func() interface{} { return new(strings.Builder) },
}

// String generates a usage page for this Meta.
//
// This method is a wrapper around the Meta.WriteMessageTo method and is untested.
func (meta Meta) String() string {
	// grab a builder from the pool
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	meta.WriteMessageTo(builder)
	return builder.String()
}

// JoinCommands joins a list of commands into a single string.
//
// This function is untested.
func JoinCommands(commands []string) string {
	// grab a builder from the pool
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	Meta{Commands: commands}.writeCommandsTo(builder)
	return builder.String()
}
