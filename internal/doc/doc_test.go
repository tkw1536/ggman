package doc_test

//spellchecker:words http httptest strings testing github cobra ggman internal
import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/doc"
)

func TestDocs_ServeHTTP(t *testing.T) {
	t.Parallel()

	// create a fake command for testing
	ggman := &cobra.Command{
		Use: "ggman",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ggman")
		},
	}

	clone := &cobra.Command{
		Use: "clone",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("clone")
		},
	}

	ls := &cobra.Command{
		Use: "ls",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ls")
		},
	}

	ggman.AddCommand(clone, ls)

	// create the docs struct
	docs, err := doc.MakeDocs(ggman)
	if err != nil {
		t.Fatalf("failed to make docs: %v", err)
	}

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{

		{
			name:           "root redirects to /ggman",
			path:           "/",
			expectedStatus: http.StatusFound,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				location := rec.Header().Get("Location")
				if location != "/ggman" {
					t.Errorf("expected redirect to /ggman, got %s", location)
				}
			},
		},
		{
			name:           "ggman command returns HTML",
			path:           "/ggman",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				contentType := rec.Header().Get("Content-Type")
				if !strings.Contains(contentType, "text/html") {
					t.Errorf("expected HTML content type, got %s", contentType)
				}

				body := rec.Body.String()
				if !strings.Contains(body, "<html") {
					t.Errorf("expected HTML content, got: %s", body)
				}
				if !strings.Contains(body, "ggman") {
					t.Errorf("expected content to mention 'ggman', got: %s", body)
				}
			},
		},

		{
			name:           "clone command returns HTML",
			path:           "/ggman/clone",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				contentType := rec.Header().Get("Content-Type")
				if !strings.Contains(contentType, "text/html") {
					t.Errorf("expected HTML content type, got %s", contentType)
				}

				body := rec.Body.String()
				if !strings.Contains(body, "<html") {
					t.Errorf("expected HTML content, got: %s", body)
				}
				if !strings.Contains(body, "clone") {
					t.Errorf("expected content to mention 'clone', got: %s", body)
				}
			},
		},

		{
			name:           "ls command returns HTML",
			path:           "/ggman/ls",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				contentType := rec.Header().Get("Content-Type")
				if !strings.Contains(contentType, "text/html") {
					t.Errorf("expected HTML content type, got %s", contentType)
				}

				body := rec.Body.String()
				if !strings.Contains(body, "<html") {
					t.Errorf("expected HTML content, got: %s", body)
				}
				if !strings.Contains(body, "ls") {
					t.Errorf("expected content to mention 'ls', got: %s", body)
				}
			},
		},
		{
			name:           "non-existent command returns 404",
			path:           "/ggman/nonexistent",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				body := rec.Body.String()
				if !strings.Contains(body, "Not Found") {
					t.Errorf("expected 'Not Found' in response, got: %s", body)
				}
			},
		},
		{
			name:           "non-existent path returns 404",
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				body := rec.Body.String()
				if !strings.Contains(body, "Not Found") {
					t.Errorf("expected 'Not Found' in response, got: %s", body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			docs.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			tt.checkResponse(t, rec)
		})
	}
}
