package doc

//spellchecker:words bytes html http strings github cobra yuin goldmark
import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/yuin/goldmark"
)

// Docs represents html documentation for a [cobra.Command].
type Docs struct {
	// Cmd is the command documentation was generated for.
	Cmd *cobra.Command

	// HTML is a map from command path to html contents.
	HTML map[string]string
}

// MakeDocs generates html documentation for a [cobra.Command].
func MakeDocs(cmd *cobra.Command) (Docs, error) {
	index := IndexFilenames(cmd, "md")

	var (
		buffer  = new(bytes.Buffer)
		builder = new(strings.Builder)
		m       = goldmark.New()
	)

	docs := Docs{
		Cmd:  cmd,
		HTML: make(map[string]string, len(index)),
	}

	if err := genDocs(cmd, &docs, buffer, builder, m, index); err != nil {
		return Docs{}, err
	}
	return docs, nil
}

func genDocs(cmd *cobra.Command, docs *Docs, buffer *bytes.Buffer, builder *strings.Builder, m goldmark.Markdown, index map[string]string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := genDocs(c, docs, buffer, builder, m, index); err != nil {
			return err
		}
	}

	buffer.Reset()
	if err := doc.GenMarkdownCustom(cmd, buffer, func(s string) string {
		return path2url(index[s])
	}); err != nil {
		return fmt.Errorf("failed to generate markdown: %w", err)
	}

	builder.Reset()

	// write a dead simple html header
	builder.WriteString("<!DOCTYPE html>")
	builder.WriteString("<html lang='en'>")
	builder.WriteString("<title>" + html.EscapeString(cmd.CommandPath()) + "</title>")
	builder.WriteString("<style>body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;}a{color:blue;}</style>")

	if err := m.Convert(buffer.Bytes(), builder); err != nil {
		return fmt.Errorf("failed to generate html: %w", err)
	}

	docs.HTML[cmd.CommandPath()] = builder.String()
	return nil
}

func path2url(path string) string {
	return "/" + strings.ReplaceAll(path, " ", "/")
}

func url2path(url string) string {
	return strings.ReplaceAll(strings.Trim(url, "/"), "/", " ")
}

// ServeHTTP implements the http.Handler interface for the Docs struct.
func (d Docs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := url2path(r.URL.Path)
	html, ok := d.HTML[path]
	if ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(html))
		return
	}
	if r.URL.Path == "/" {
		http.Redirect(w, r, path2url(d.Cmd.Root().CommandPath()), http.StatusFound)
		return
	}

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
