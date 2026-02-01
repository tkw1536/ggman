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

// MakeDocs generates html documentation for a [cobra.Command].
func MakeDocs(cmd *cobra.Command) (Docs, error) {
	var g generator
	if err := g.build(cmd); err != nil {
		return Docs{}, err
	}

	return Docs{
		rootCommandPath: cmd.CommandPath(),
		html:            g.html,
	}, nil
}

type Docs struct {
	rootCommandPath string
	html            map[string]string
}

// ServeHTTP implements the http.Handler interface for the Docs struct.
func (d Docs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := url2path(r.URL.Path)
	html, ok := d.html[path]
	if ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(html))
		return
	}
	if r.URL.Path == "/" {
		http.Redirect(w, r, path2url(d.rootCommandPath), http.StatusFound)
		return
	}

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

type generator struct {
	// index holds a map from filename to name of a command.
	index map[string]string

	buffer  bytes.Buffer
	builder strings.Builder
	m       goldmark.Markdown

	// html holds a map from command path to html contents.
	html map[string]string
}

// build builds documentation for a command.
func (g *generator) build(root *cobra.Command) error {
	g.index = IndexFilenames(root, "md")

	g.m = goldmark.New()
	g.html = make(map[string]string, len(g.index))

	if err := g.genDocs(root); err != nil {
		return fmt.Errorf("failed to generate docs: %w", err)
	}
	return nil
}

func (g *generator) genDocs(cmd *cobra.Command) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := g.genDocs(c); err != nil {
			return err
		}
	}

	g.buffer.Reset()
	if err := doc.GenMarkdownCustom(cmd, &g.buffer, func(s string) string {
		return path2url(g.index[s])
	}); err != nil {
		return fmt.Errorf("failed to generate markdown: %w", err)
	}

	g.builder.Reset()

	// write a dead simple html header
	g.builder.WriteString("<!DOCTYPE html>")
	g.builder.WriteString("<html lang='en'>")
	g.builder.WriteString("<title>" + html.EscapeString(cmd.CommandPath()) + "</title>")
	g.builder.WriteString("<style>body{font-family:-apple-system,BlinkMacSystemFont,Helvetica,Arial,sans-serif;}a{color:blue;}</style>")

	if err := g.m.Convert(g.buffer.Bytes(), &g.builder); err != nil {
		return fmt.Errorf("failed to generate html: %w", err)
	}

	g.html[cmd.CommandPath()] = g.builder.String()
	return nil
}

func path2url(path string) string {
	return "/" + strings.ReplaceAll(path, " ", "/")
}

func url2path(url string) string {
	return strings.ReplaceAll(strings.Trim(url, "/"), "/", " ")
}
