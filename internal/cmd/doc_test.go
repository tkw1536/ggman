package cmd_test

//spellchecker:words context http strconv testing ggman internal mockenv testutil
import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"testing"

	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
	"go.tkw01536.de/ggman/internal/testutil"
)

func TestCommandDoc(t *testing.T) {
	t.Parallel()

	var (
		port = strconv.Itoa(testutil.FindFreePort())
		host = "127.0.0.1"
		addr = net.JoinHostPort(host, port)
	)

	// Prepare to run the doc command with the chosen port and no browser open.
	args := []string{"doc", "--host", host, "--port", port, "--no-open"}

	ctx, cancel := context.WithCancel(t.Context())

	// start a goroutine to wait for and check if the server is listening
	errs := make(chan string, 1)
	go func() {
		defer cancel()

		if err := testutil.WaitForPort(t.Context(), addr); err != nil {
			errs <- fmt.Sprintf("failed to wait for port: %v", err)
			return
		}

		// create a new client that does not follow redirects
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://"+addr, nil)
		if err != nil {
			errs <- fmt.Sprintf("failed to create request: %v", err)
			return
		}
		res, err := client.Do(req)
		if err != nil {
			errs <- fmt.Sprintf("failed to get: %v", err)
			return
		}
		defer res.Body.Close() //nolint:errcheck // we don't care about a failed close in a test
		if res.StatusCode != http.StatusFound {
			errs <- fmt.Sprintf("expected http.StatusFound, got %d", res.StatusCode)
			return
		}

		if location := res.Header.Get("Location"); location != "/ggman" {
			errs <- fmt.Sprintf("expected Location header %q, got %q", "/ggman", location)
			return
		}
	}()

	mock := mockenv.NewMockEnv(t)
	code, stdout, stderr := mock.Run(t, ctx, cmd.NewCommand, "", "", args...)
	if code != 0 {
		t.Errorf("command failed: %d", code)
	}

	// check that the output is as expected
	mock.AssertOutput(t, "Stdout", stdout, "server listening at http://"+addr+"\n")
	mock.AssertOutput(t, "Stderr", stderr, "")

	select {
	case err := <-errs:
		t.Error(err)
	default:
	}
}
