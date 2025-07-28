//spellchecker:words testutil
package testutil

//spellchecker:words context time
import (
	"context"
	"fmt"
	"net"
	"time"
)

// FindFreePort picks a random non-zero unassigned port on localhost.
// It is only guaranteed to be free at the time the function is invoked, and other programs may claim it.
// If no free port is found, the function panics.
func FindFreePort() int {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Errorf("failed to find free port: %w", err))
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		panic(fmt.Errorf("failed to close listener: %w", err))
	}
	if port == 0 {
		panic("free port is 0")
	}
	return port
}

const waitPortInterval = 10 * time.Millisecond

// WaitForPort waits until the given address is reachable via TCP.
func WaitForPort(ctx context.Context, addr string) error {
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled: %w", err)
		}

		// try to dial and close if successful
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			if err := conn.Close(); err != nil {
				return fmt.Errorf("failed to close connection: %w", err)
			}
			return nil
		}

		// wait a bit and try again
		time.Sleep(waitPortInterval)
	}
}
