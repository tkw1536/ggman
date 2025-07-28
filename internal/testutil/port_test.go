//spellchecker:words testutil
package testutil_test

//spellchecker:words context strconv ggman internal testutil
import (
	"context"
	"fmt"
	"net"
	"strconv"

	"go.tkw01536.de/ggman/internal/testutil"
)

// ExampleFindFreePort demonstrates how to use FindFreePort.
func ExampleFindFreePort() {
	port := testutil.FindFreePort()
	if port == 0 {
		panic("free port is 0")
	}
	fmt.Println("picked a free port")

	// Output:
	// picked a free port
}

func ExampleWaitForPort() {
	// pick a random port
	port := testutil.FindFreePort()
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))

	waitPortReturned := make(chan struct{})
	// wait for that port do be available
	go func() {
		defer close(waitPortReturned)

		err := testutil.WaitForPort(context.Background(), addr)
		fmt.Printf("WaitForPort returned: %v\n", err)
	}()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	<-waitPortReturned
	if err := listener.Close(); err != nil {
		panic(err)
	}
	fmt.Println("listener closed")

	// Output:
	// WaitForPort returned: <nil>
	// listener closed
}
