package cmd

//spellchecker:words github cobra ggman constants legal
import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/doc"
	"go.tkw01536.de/pkglib/exit"
)

var (
	errDocGenerate    = exit.NewErrorWithCode("failed to generate docs", exit.ExitGeneric)
	errServerListen   = exit.NewErrorWithCode("failed to listen", exit.ExitGeneric)
	errServerShutdown = exit.NewErrorWithCode("failed to shutdown server", exit.ExitGeneric)
	errDocOpenBrowser = exit.NewErrorWithCode("failed to open browser", exit.ExitGeneric)
)

//spellchecker:words nolint wrapcheck

func NewDocCommand() *cobra.Command {
	impl := new(_doc)

	cmd := &cobra.Command{
		Use:   "doc",
		Short: "start a server with ggman documentation",
		Long:  "",
		Args:  cobra.NoArgs,

		RunE: impl.Exec,
	}

	flags := cmd.Flags()
	flags.StringVarP(&impl.Host, "host", "", "localhost", "host to listen on")
	flags.IntVarP(&impl.Port, "port", "p", 0, "port to listen on")
	flags.BoolVarP(&impl.NoOpen, "no-open", "", false, "don't open the browser")

	return cmd
}

type _doc struct {
	Host   string
	Port   int
	NoOpen bool
}

func (d *_doc) Exec(cmd *cobra.Command, args []string) (e error) {
	docs, err := doc.MakeDocs(cmd.Root())
	if err != nil {
		return fmt.Errorf("%w: %w", errDocGenerate, err)
	}

	server := &http.Server{
		Handler:           docs,
		ReadHeaderTimeout: 30 * time.Second,
	}

	// start listening for connections
	l, err := net.Listen("tcp", net.JoinHostPort(d.Host, strconv.Itoa(d.Port)))
	if err != nil {
		return fmt.Errorf("%w: %w", errServerListen, err)
	}
	defer l.Close() //nolint:errcheck // we don't care about a failed close

	// determine the address of the server
	addr := "http://" + l.Addr().String()
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "server listening at %s\n", addr); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}

	// open in browser if requested
	if !d.NoOpen {
		if err := browser.OpenURL(addr); err != nil {
			return fmt.Errorf("%w: %w", errDocOpenBrowser, err)
		}
	}

	errChan := make(chan error, 1)

	//nolint:errcheck // sending the error to errChan
	go func() (err error) {
		defer func() {
			if err == nil {
				return
			}
			errChan <- err
		}()

		if err := server.Serve(l); err != nil {
			return fmt.Errorf("%w: %w", errServerListen, err)
		}
		return nil
	}()

	select {
	case <-cmd.Context().Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("%w: %w", errServerShutdown, err)
		}
		return nil
	case err := <-errChan:
		return err
	}
}
