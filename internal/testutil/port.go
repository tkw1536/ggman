//spellchecker:words testutil
package testutil

//spellchecker:words context time
import (
	"context"
	"fmt"
	"net"
	"time"
)

// FindFreePort picks a random non-zero unassigned port on the given ost.
// It is only guaranteed to be free at the time the function is invoked, and other programs may claim it.
// If no free port is found, or ctx expires, the function panics.
func FindFreePort(ctx context.Context, host string) int {
	var lc net.ListenConfig
	l, err := lc.Listen(ctx, "tcp", host+":0")
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
	var dialer net.Dialer

	timer := time.NewTimer(waitPortInterval)
	for {
		// try to dial and close if successful
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			if err := conn.Close(); err != nil {
				return fmt.Errorf("failed to close connection: %w", err)
			}
			return nil
		}

		// wait to try again, or close if the context is done.
		timer.Reset(waitPortInterval)
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		case <-timer.C:
		}
	}
}
