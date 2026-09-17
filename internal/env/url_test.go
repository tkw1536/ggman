package env_test

//spellchecker:words reflect testing ggman internal
import (
	"reflect"
	"testing"

	"go.tkw01536.de/ggman/internal/env"
)

var urlTests = []struct {
	name string
	str  string
	url  env.URL
}{
	// ssh://[user@]host.xz[:port]/path/to/repo.git/
	{
		name: "ssh",
		str:  "ssh://host.xz/path/to/repo.git/",
		url:  env.URL{Scheme: "ssh", User: "", Password: "", HostName: "host.xz", Port: 0, Path: "path/to/repo.git/"},
	},
	{
		name: "sshUser",
		str:  "ssh://user@host.xz/path/to/repo.git/",
		url:  env.URL{Scheme: "ssh", User: "user", Password: "", HostName: "host.xz", Port: 0, Path: "path/to/repo.git/"},
	},
	{
		name: "sshPort",
		str:  "ssh://host.xz:1234/path/to/repo.git/",
		url:  env.URL{Scheme: "ssh", User: "", Password: "", HostName: "host.xz", Port: 1234, Path: "path/to/repo.git/"},
	},
	{
		name: "sshUserPort",
		str:  "ssh://user@host.xz:1234/path/to/repo.git/",
		url:  env.URL{Scheme: "ssh", User: "user", Password: "", HostName: "host.xz", Port: 1234, Path: "path/to/repo.git/"},
	},

	// git://host.xz[:port]/path/to/repo.git/
	{
		name: "git",
		str:  "git://host.xz/path/to/repo.git/",
		url:  env.URL{Scheme: "git", User: "", Password: "", HostName: "host.xz", Port: 0, Path: "path/to/repo.git/"},
	},

	{
		name: "gitPort",
		str:  "git://host.xz:1234/path/to/repo.git/",
		url:  env.URL{Scheme: "git", User: "", Password: "", HostName: "host.xz", Port: 1234, Path: "path/to/repo.git/"},
	},

	//  [user@]host.xz:path/to/repo.git/
	{
		name: "noProto",
		str:  "host.xz:path/to/repo.git/",
		url:  env.URL{Scheme: "", User: "", Password: "", HostName: "host.xz", Port: 0, Path: "path/to/repo.git/"},
	},
	{
		name: "noProtoUser",
		str:  "user@host.xz:path/to/repo.git/",
		url:  env.URL{Scheme: "", User: "user", Password: "", HostName: "host.xz", Port: 0, Path: "path/to/repo.git/"},
	},

	// local paths
	{
		name: "localFile",
		str:  "file:///path/to/somewhere",
		url:  env.URL{Scheme: "file", User: "", Password: "", HostName: "", Port: 0, Path: "path/to/somewhere"},
	},

	{
		name: "localPath",
		str:  "/path/to/somewhere",
		url:  env.URL{Scheme: "", User: "", Password: "", HostName: "", Port: 0, Path: "path/to/somewhere"},
	},

	{
		name: "localRelPath",
		str:  "../some/relative/path",
		url:  env.URL{Scheme: "", User: "", Password: "", HostName: "..", Port: 0, Path: "some/relative/path"},
	},

	{
		name: "localRelPath2",
		str:  "./some/relative/path",
		url:  env.URL{Scheme: "", User: "", Password: "", HostName: ".", Port: 0, Path: "some/relative/path"},
	},
}

func TestParseURL(t *testing.T) {
	t.Parallel()

	for _, tt := range urlTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotURL := env.ParseURL(tt.str)
			if !reflect.DeepEqual(gotURL, tt.url) {
				t.Errorf("env.ParseURL() = %v, want %v", gotURL, tt.url)
			}
		})
	}
}

func TestURL_SplitForgeReference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		url          env.URL
		wantPath     string
		wantSuffixes string
		wantRef      string
		wantRelative string
	}{
		{
			name:         "tree with ref and relative",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo/tree/ref/path"},
			wantPath:     "user/repo",
			wantRef:      "ref",
			wantRelative: "path",
		},
		{
			name:         "tree with ref only",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo/tree/main"},
			wantPath:     "user/repo",
			wantRef:      "main",
			wantRelative: "",
		},
		{
			name:         "tree with ref only trailing slash",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo/tree/main/"},
			wantPath:     "user/repo",
			wantRef:      "main",
			wantRelative: "",
		},
		{
			name:         "tree with nested relative path",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo/tree/main/internal/env"},
			wantPath:     "user/repo",
			wantRef:      "main",
			wantRelative: "internal/env",
		},
		{
			name:         "optional /_/ before tree",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo/_/tree/main/src"},
			wantPath:     "user/repo",
			wantRef:      "main",
			wantRelative: "src",
		},
		{
			name:         "query suffix",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo/tree/main/src?tab=readme"},
			wantPath:     "user/repo",
			wantSuffixes: "?tab=readme",
			wantRef:      "main",
			wantRelative: "src",
		},
		{
			name:         "fragment suffix",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo/tree/main/src#L10"},
			wantPath:     "user/repo",
			wantSuffixes: "#L10",
			wantRef:      "main",
			wantRelative: "src",
		},
		{
			name:         "query then fragment",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo/tree/main/src?hello=world#fragment"},
			wantPath:     "user/repo",
			wantSuffixes: "?hello=world#fragment",
			wantRef:      "main",
			wantRelative: "src",
		},
		{
			name:         "fragment containing question mark",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo/tree/main/src#fragment?also=fragment"},
			wantPath:     "user/repo",
			wantSuffixes: "#fragment?also=fragment",
			wantRef:      "main",
			wantRelative: "src",
		},
		{
			name:         "no tree reference",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo"},
			wantPath:     "user/repo",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "no tree reference with query",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo?foo=bar"},
			wantPath:     "user/repo",
			wantSuffixes: "?foo=bar",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "no tree reference with query then fragment",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "some/path/?hello=world#fragment"},
			wantPath:     "some/path/",
			wantSuffixes: "?hello=world#fragment",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "no tree reference with fragment containing question mark",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "some/path#fragment?also=fragment"},
			wantPath:     "some/path",
			wantSuffixes: "#fragment?also=fragment",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "trailing _ in path without tree reference",
			url:          env.URL{Scheme: "https", HostName: "example.com", Path: "some/path/_"},
			wantPath:     "some/path/_",
			wantSuffixes: "",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "scheme-less host is web url",
			url:          env.URL{HostName: "example.com", Path: "user/repo/tree/main/src"},
			wantPath:     "user/repo",
			wantRef:      "main",
			wantRelative: "src",
		},
		{
			name:         "ssh with tree is not a web url",
			url:          env.URL{Scheme: "ssh", User: "git", HostName: "gitforge.example", Path: "user/repo/tree/main/src"},
			wantPath:     "user/repo/tree/main/src",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "ssh with query and fragment is not a web url",
			url:          env.URL{Scheme: "ssh", HostName: "gitforge.example", Path: "user/repo/tree/main/src?tab=readme#L10"},
			wantPath:     "user/repo/tree/main/src?tab=readme#L10",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "git with tree is not a web url",
			url:          env.URL{Scheme: "git", HostName: "host.xz", Path: "path/to/repo/tree/main"},
			wantPath:     "path/to/repo/tree/main",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "file with tree is not a web url",
			url:          env.URL{Scheme: "file", Path: "path/to/repo/tree/main/src"},
			wantPath:     "path/to/repo/tree/main/src",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "local path with tree is not a web url",
			url:          env.URL{Path: "user/repo/tree/main/src?foo=bar"},
			wantPath:     "user/repo/tree/main/src?foo=bar",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "gitlab-style /- /tree",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/-/tree/my-branch/path"},
			wantPath:     "user/repo",
			wantRef:      "my-branch",
			wantRelative: "path",
		},
		{
			name:         "gitlab-style /- /tree ref only",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/-/tree/my-branch"},
			wantPath:     "user/repo",
			wantRef:      "my-branch",
			wantRelative: "",
		},
		{
			name:         "gitlab-style /- /tree ref only trailing slash",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/-/tree/my-branch/"},
			wantPath:     "user/repo",
			wantRef:      "my-branch",
			wantRelative: "",
		},
		{
			name:         "gitea src/branch with path",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/src/branch/my-branch/path"},
			wantPath:     "user/repo",
			wantRef:      "my-branch",
			wantRelative: "path",
		},
		{
			name:         "gitea src/tag with path",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/src/tag/v1.0/path"},
			wantPath:     "user/repo",
			wantRef:      "v1.0",
			wantRelative: "path",
		},
		{
			name:         "gitea src/branch root of repo",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/src/branch/main"},
			wantPath:     "user/repo",
			wantRef:      "main",
			wantRelative: "",
		},
		{
			name:         "gitea src/branch root of repo trailing slash",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/src/branch/main/"},
			wantPath:     "user/repo",
			wantRef:      "main",
			wantRelative: "",
		},
		{
			name:         "gitea src/tag root of repo",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/src/tag/v1.0"},
			wantPath:     "user/repo",
			wantRef:      "v1.0",
			wantRelative: "",
		},
		{
			name:         "gitea src/tag root of repo trailing slash",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/src/tag/v1.0/"},
			wantPath:     "user/repo",
			wantRef:      "v1.0",
			wantRelative: "",
		},
		{
			name:         "sourcehut-style tree",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "~user/repo/tree/my-branch/path"},
			wantPath:     "~user/repo",
			wantRef:      "my-branch",
			wantRelative: "path",
		},
		{
			name:         "sourcehut-style tree ref only",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "~user/repo/tree/my-branch"},
			wantPath:     "~user/repo",
			wantRef:      "my-branch",
			wantRelative: "",
		},
		{
			name:         "sourcehut-style tree ref only trailing slash",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "~user/repo/tree/my-branch/"},
			wantPath:     "~user/repo",
			wantRef:      "my-branch",
			wantRelative: "",
		},
		{
			name:         "blob with ref and relative",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/blob/main/README.md"},
			wantPath:     "user/repo",
			wantRef:      "main",
			wantRelative: "README.md",
		},
		{
			name:         "blob ref only",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/blob/main"},
			wantPath:     "user/repo",
			wantRef:      "main",
			wantRelative: "",
		},
		{
			name:         "blob ref only trailing slash",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/blob/main/"},
			wantPath:     "user/repo",
			wantRef:      "main",
			wantRelative: "",
		},
		{
			name:         "encoded slash in ref",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/tree/feature%2Fx/path"},
			wantPath:     "user/repo",
			wantRef:      "feature%2Fx",
			wantRelative: "path",
		},
		{
			name:         "unencoded slash splits ref and relative",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/tree/feature/x"},
			wantPath:     "user/repo",
			wantRef:      "feature",
			wantRelative: "x",
		},
		{
			name:         "no false positive restructure",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/restructure/docs"},
			wantPath:     "user/restructure/docs",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "no false positive treehouse",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/treehouse/docs"},
			wantPath:     "user/treehouse/docs",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "src alone is not a marker",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo/src/docs"},
			wantPath:     "user/repo/src/docs",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "top-level tree/* repository is not a forge reference",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "tree/main"},
			wantPath:     "tree/main",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "top-level tree/* repository with query",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "tree/main?tab=readme"},
			wantPath:     "tree/main",
			wantSuffixes: "?tab=readme",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "top-level blob/* repository is not a forge reference",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "blob/main"},
			wantPath:     "blob/main",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "top-level src/* repository is not a forge reference",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "src/branch"},
			wantPath:     "src/branch",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "single-segment tree path is not a forge reference",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "tree"},
			wantPath:     "tree",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "top-level src/* repository is not a forge reference",
			url:          env.URL{Scheme: "https", HostName: "gitforge.example", Path: "src/branch/example"},
			wantPath:     "src/branch/example",
			wantRef:      "",
			wantRelative: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotPath, gotSuffixes, gotRef, gotRelative := tt.url.SplitForgeReference()
			if gotPath != tt.wantPath || gotSuffixes != tt.wantSuffixes || gotRef != tt.wantRef || gotRelative != tt.wantRelative {
				t.Errorf("URL.SplitForgeReference() = (%q, %q, %q, %q), want (%q, %q, %q, %q)",
					gotPath, gotSuffixes, gotRef, gotRelative,
					tt.wantPath, tt.wantSuffixes, tt.wantRef, tt.wantRelative)
			}
		})
	}
}

func TestParseURLAndForgeReference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        string
		wantURL      env.URL
		wantSuffix   string
		wantRef      string
		wantRelative string
	}{
		{
			name:         "tree with ref and relative",
			input:        "https://example.com/user/repo/tree/ref/path",
			wantURL:      env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo"},
			wantRef:      "ref",
			wantRelative: "path",
		},
		{
			name:         "tree with ref only",
			input:        "https://example.com/user/repo/tree/main",
			wantURL:      env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo"},
			wantRef:      "main",
			wantRelative: "",
		},
		{
			name:         "tree with nested relative path",
			input:        "https://example.com/user/repo/tree/main/internal/env",
			wantURL:      env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo"},
			wantRef:      "main",
			wantRelative: "internal/env",
		},
		{
			name:         "optional /_/ before tree",
			input:        "https://example.com/user/repo/_/tree/main/src",
			wantURL:      env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo"},
			wantRef:      "main",
			wantRelative: "src",
		},
		{
			name:         "query suffix",
			input:        "https://example.com/user/repo/tree/main/src?tab=readme",
			wantURL:      env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo"},
			wantSuffix:   "?tab=readme",
			wantRef:      "main",
			wantRelative: "src",
		},
		{
			name:         "fragment suffix",
			input:        "https://example.com/user/repo/tree/main/src#L10",
			wantURL:      env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo"},
			wantSuffix:   "#L10",
			wantRef:      "main",
			wantRelative: "src",
		},
		{
			name:         "query then fragment",
			input:        "https://example.com/user/repo/tree/main/src?hello=world#fragment",
			wantURL:      env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo"},
			wantSuffix:   "?hello=world#fragment",
			wantRef:      "main",
			wantRelative: "src",
		},
		{
			name:         "fragment containing question mark",
			input:        "https://example.com/user/repo/tree/main/src#fragment?also=fragment",
			wantURL:      env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo"},
			wantSuffix:   "#fragment?also=fragment",
			wantRef:      "main",
			wantRelative: "src",
		},
		{
			name:         "no tree reference",
			input:        "https://example.com/user/repo",
			wantURL:      env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo"},
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "no tree reference with query then fragment",
			input:        "https://example.com/some/path/?hello=world#fragment",
			wantURL:      env.URL{Scheme: "https", HostName: "example.com", Path: "some/path/"},
			wantSuffix:   "?hello=world#fragment",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "no tree reference with fragment containing question mark",
			input:        "https://example.com/some/path#fragment?also=fragment",
			wantURL:      env.URL{Scheme: "https", HostName: "example.com", Path: "some/path"},
			wantSuffix:   "#fragment?also=fragment",
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "scheme-less host with tree",
			input:        "example.com:user/repo/tree/main/src",
			wantURL:      env.URL{HostName: "example.com", Path: "user/repo"},
			wantRef:      "main",
			wantRelative: "src",
		},
		{
			name:         "ssh url with tree is not a web url",
			input:        "ssh://git@gitforge.example/user/repo/tree/main/src",
			wantURL:      env.URL{Scheme: "ssh", User: "git", HostName: "gitforge.example", Path: "user/repo/tree/main/src"},
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "ssh url with query and fragment is not a web url",
			input:        "ssh://git@gitforge.example/user/repo/tree/main/src?tab=readme#L10",
			wantURL:      env.URL{Scheme: "ssh", User: "git", HostName: "gitforge.example", Path: "user/repo/tree/main/src?tab=readme#L10"},
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "git url with tree is not a web url",
			input:        "git://host.xz/path/to/repo/tree/main",
			wantURL:      env.URL{Scheme: "git", HostName: "host.xz", Path: "path/to/repo/tree/main"},
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "file url with tree is not a web url",
			input:        "file:///path/to/repo/tree/main/src",
			wantURL:      env.URL{Scheme: "file", Path: "path/to/repo/tree/main/src"},
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "local path with tree is not a web url",
			input:        "/user/repo/tree/main/src?foo=bar",
			wantURL:      env.URL{Path: "user/repo/tree/main/src?foo=bar"},
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "gitlab-style /- /tree",
			input:        "https://gitforge.example/user/repo/-/tree/my-branch/path",
			wantURL:      env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo"},
			wantRef:      "my-branch",
			wantRelative: "path",
		},
		{
			name:         "gitea src/branch",
			input:        "https://gitforge.example/user/repo/src/branch/my-branch/path",
			wantURL:      env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo"},
			wantRef:      "my-branch",
			wantRelative: "path",
		},
		{
			name:         "gitea src/tag",
			input:        "https://gitforge.example/user/repo/src/tag/v1.0",
			wantURL:      env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo"},
			wantRef:      "v1.0",
			wantRelative: "",
		},
		{
			name:         "sourcehut-style tree",
			input:        "https://gitforge.example/~user/repo/tree/my-branch/path",
			wantURL:      env.URL{Scheme: "https", HostName: "gitforge.example", Path: "~user/repo"},
			wantRef:      "my-branch",
			wantRelative: "path",
		},
		{
			name:         "blob with path",
			input:        "https://gitforge.example/user/repo/blob/main/README.md",
			wantURL:      env.URL{Scheme: "https", HostName: "gitforge.example", Path: "user/repo"},
			wantRef:      "main",
			wantRelative: "README.md",
		},
		{
			name:         "top-level tree/* repository is not a forge reference",
			input:        "https://gitforge.example/tree/main",
			wantURL:      env.URL{Scheme: "https", HostName: "gitforge.example", Path: "tree/main"},
			wantRef:      "",
			wantRelative: "",
		},
		{
			name:         "top-level src/* repository is not a forge reference",
			input:        "https://gitforge.example/src/branch",
			wantURL:      env.URL{Scheme: "https", HostName: "gitforge.example", Path: "src/branch"},
			wantRef:      "",
			wantRelative: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotURL, gotSuffix, gotRef, gotRelative := env.ParseURLAndForgeReference(tt.input)
			if !reflect.DeepEqual(gotURL, tt.wantURL) || gotSuffix != tt.wantSuffix || gotRef != tt.wantRef || gotRelative != tt.wantRelative {
				t.Errorf("ParseURLAndForgeReference(%q) = (%v, %q, %q, %q), want (%v, %q, %q, %q)",
					tt.input, gotURL, gotSuffix, gotRef, gotRelative,
					tt.wantURL, tt.wantSuffix, tt.wantRef, tt.wantRelative)
			}
		})
	}
}

func TestURL_String(t *testing.T) {
	t.Parallel()

	for _, tt := range urlTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotString := tt.url.String()
			if gotString != tt.str {
				t.Errorf("URL.String() = %v, want %v", gotString, tt.str)
			}
		})
	}
}

func Benchmark_ParseURL(b *testing.B) {
	for b.Loop() {
		env.ParseURL("ssh://host.xz/path/to/repo.git/")
		env.ParseURL("ssh://user@host.xz/path/to/repo.git/")
		env.ParseURL("ssh://host.xz:1234/path/to/repo.git/")
		env.ParseURL("ssh://user@host.xz:1234/path/to/repo.git/")
		env.ParseURL("git://host.xz/path/to/repo.git/")
		env.ParseURL("git://host.xz:1234/path/to/repo.git/")
		env.ParseURL("host.xz:path/to/repo.git/")
		env.ParseURL("user@host.xz:path/to/repo.git/")
	}
}

var urlLocalityTests = []struct {
	name     string
	url      env.URL
	isLocal  bool
	isWebURL bool
}{
	{
		"ssh",
		env.URL{"ssh", "", "", "host.xz", 0, "path/to/repo.git/"},
		false,
		false,
	},
	{
		"sshUser",
		env.URL{"ssh", "user", "", "host.xz", 0, "path/to/repo.git/"},
		false,
		false,
	},
	{
		"sshPort",
		env.URL{"ssh", "", "", "host.xz", 1234, "path/to/repo.git/"},
		false,
		false,
	},
	{
		"sshUserPort",
		env.URL{"ssh", "user", "", "host.xz", 1234, "path/to/repo.git/"},
		false,
		false,
	},

	// git://host.xz[:port]/path/to/repo.git/
	{
		"git",
		env.URL{"git", "", "", "host.xz", 0, "path/to/repo.git/"},
		false,
		false,
	},
	{
		"gitPort",
		env.URL{"git", "", "", "host.xz", 1234, "path/to/repo.git/"},
		false,
		false,
	},

	//  [user@]host.xz:path/to/repo.git/
	{
		"noProto",
		env.URL{"", "", "", "host.xz", 0, "path/to/repo.git/"},
		false,
		true,
	},
	{
		"noProtoUser",
		env.URL{"", "user", "", "host.xz", 0, "path/to/repo.git/"},
		false,
		true,
	},

	// http(s) web urls
	{
		"https",
		env.URL{"https", "", "", "example.com", 0, "user/repo"},
		false,
		true,
	},
	{
		"http",
		env.URL{"http", "", "", "example.com", 0, "user/repo"},
		false,
		true,
	},
	{
		"httpsUpper",
		env.URL{"HTTPS", "", "", "example.com", 0, "user/repo"},
		false,
		true,
	},
	{
		"httpUpper",
		env.URL{"HTTP", "", "", "example.com", 0, "user/repo"},
		false,
		true,
	},
	{
		"httpsUserPort",
		env.URL{"https", "user", "", "example.com", 8443, "user/repo"},
		false,
		true,
	},

	// local paths
	{
		"localFile",
		env.URL{"file", "", "", "", 0, "path/to/somewhere"},
		true,
		false,
	},
	{
		"localPath",
		env.URL{"", "", "", "", 0, "path/to/somewhere"},
		true,
		false,
	},
	{
		"localRelPath",
		env.URL{"", "", "", "..", 0, "some/relative/path"},
		true,
		false,
	},
	{
		"localRelPath2",
		env.URL{"", "", "", ".", 0, "some/relative/path"},
		true,
		false,
	},
}

func TestURL_IsLocal(t *testing.T) {
	t.Parallel()

	for _, tt := range urlLocalityTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.url.IsLocal(); got != tt.isLocal {
				t.Errorf("URL.IsLocal() = %v, want %v", got, tt.isLocal)
			}
		})
	}
}

func TestURL_IsWebURL(t *testing.T) {
	t.Parallel()

	for _, tt := range urlLocalityTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.url.IsWebURL(); got != tt.isWebURL {
				t.Errorf("URL.IsWebURL() = %v, want %v", got, tt.isWebURL)
			}
		})
	}
}

func TestURL_IsRelative(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url  env.URL
		want bool
	}{
		{"ssh remote", env.URL{Scheme: "ssh", HostName: "host.xz", Path: "path/to/repo.git/"}, false},
		{"https remote", env.URL{Scheme: "https", HostName: "example.com", Path: "user/repo"}, false},
		{"no scheme remote", env.URL{HostName: "host.xz", Path: "path/to/repo.git/"}, false},
		{"empty url", env.URL{}, false},
		{"absolute local path", env.URL{Path: "path/to/somewhere"}, false},
		{"file url", env.URL{Scheme: "file", Path: "path/to/somewhere"}, false},
		{"hidden git dir", env.URL{HostName: "host.xz", Path: "a/.git/b"}, false},

		{"hostname is dot", env.URL{HostName: ".", Path: "some/relative/path"}, true},
		{"hostname is parent", env.URL{HostName: "..", Path: "some/relative/path"}, true},
		{"user is dot", env.URL{User: ".", HostName: "host.xz", Path: "repo"}, true},
		{"user is parent", env.URL{User: "..", HostName: "host.xz", Path: "repo"}, true},
		{"git user is ignored", env.URL{User: "git", HostName: "host.xz", Path: "repo"}, false},

		{"path is parent", env.URL{HostName: "host.xz", Path: ".."}, true},
		{"path is dot", env.URL{HostName: "host.xz", Path: "."}, false},
		{"leading parent in path", env.URL{HostName: "host.xz", Path: "../repo"}, true},
		{"multiple leading parents in path", env.URL{HostName: "host.xz", Path: "../../repo"}, true},
		{"parent beyond root of path", env.URL{HostName: "host.xz", Path: "a/../.."}, true},
		{"parent then file beyond root", env.URL{HostName: "host.xz", Path: "a/../../b"}, true},

		{"dot segments are stripped", env.URL{HostName: "host.xz", Path: "a/./b"}, false},
		{"resolved parent in path", env.URL{HostName: "host.xz", Path: "a/../b"}, false},
		{"resolved chained parents", env.URL{HostName: "host.xz", Path: "a/b/../.."}, false},
		{"mixed dots and resolved parents", env.URL{HostName: "host.xz", Path: "a/./b/../c"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.url.IsRelative(); got != tt.want {
				t.Errorf("URL.IsRelative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURL_PathSegments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url  env.URL
		want []string
	}{
		{"empty path", env.URL{Path: ""}, []string{}},
		{"single segment", env.URL{Path: "file"}, []string{"file"}},
		{"relative path", env.URL{Path: "a/b/c"}, []string{"a", "b", "c"}},
		{"keeps .git suffix", env.URL{Path: "hello/world.git"}, []string{"hello", "world.git"}},

		{"leading slash", env.URL{Path: "/a/b/c"}, []string{"a", "b", "c"}},
		{"trailing slash", env.URL{Path: "a/b/c/"}, []string{"a", "b", "c"}},
		{"leading and trailing slash", env.URL{Path: "/a/b/c/"}, []string{"a", "b", "c"}},
		{"only slashes", env.URL{Path: "/////"}, []string{}},
		{"repeated slashes", env.URL{Path: "a//b///c"}, []string{"a", "b", "c"}},
		{"repeated slashes at start and end", env.URL{Path: "//a/b/c//"}, []string{"a", "b", "c"}},

		{"dot", env.URL{Path: "."}, []string{}},
		{"dot slash", env.URL{Path: "./"}, []string{}},
		{"slash dot", env.URL{Path: "/."}, []string{}},
		{"dot segment", env.URL{Path: "a/./b"}, []string{"a", "b"}},
		{"leading dot segment", env.URL{Path: "./a/b"}, []string{"a", "b"}},
		{"trailing dot segment", env.URL{Path: "a/b/."}, []string{"a", "b"}},
		{"only dots", env.URL{Path: "././."}, []string{}},
		{"hidden file is kept", env.URL{Path: "a/.git/b"}, []string{"a", ".git", "b"}},
		{"dotfile", env.URL{Path: ".gitignore"}, []string{".gitignore"}},

		{"parent", env.URL{Path: ".."}, []string{".."}},
		{"slash parent", env.URL{Path: "/.."}, []string{".."}},
		{"parent of file", env.URL{Path: "a/b/.."}, []string{"a"}},
		{"parent in the middle", env.URL{Path: "a/b/../c"}, []string{"a", "c"}},
		{"leading parent is kept", env.URL{Path: "../a/b"}, []string{"..", "a", "b"}},
		{"multiple leading parents are kept", env.URL{Path: "../../a"}, []string{"..", "..", "a"}},
		{"parent beyond root of relative path", env.URL{Path: "a/../.."}, []string{".."}},
		{"parent then file beyond root", env.URL{Path: "a/../../b"}, []string{"..", "b"}},
		{"chained parents", env.URL{Path: "a/b/c/../../d"}, []string{"a", "d"}},
		{"all parents of relative path", env.URL{Path: "a/b/../.."}, []string{}},
		{"all parents of absolute path", env.URL{Path: "/a/b/../.."}, []string{}},
		{"parent does not remove another parent", env.URL{Path: "../../../a"}, []string{"..", "..", "..", "a"}},

		{"mixed dots and parents", env.URL{Path: "a/./b/../c"}, []string{"a", "c"}},
		{"mixed slashes, dots and parents", env.URL{Path: "/a//./b/../c//"}, []string{"a", "c"}},
		{"dot then parent", env.URL{Path: "a/./../b"}, []string{"b"}},
		{"parent then dot", env.URL{Path: "a/.././b"}, []string{"b"}},
		{"current dir then parent", env.URL{Path: "./.."}, []string{".."}},
		{"parent then current dir", env.URL{Path: "../."}, []string{".."}},

		{
			"ignores other URL fields",
			env.URL{Scheme: "ssh", User: "git", Password: "secret", HostName: "host.xz", Port: 22, Path: "hello/world"},
			[]string{"hello", "world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.url.PathSegments(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("URL.PathSegments() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestURL_Components(t *testing.T) {
	t.Parallel()

	type fields struct {
		Scheme   string
		User     string
		Password string
		HostName string
		Port     uint16
		Path     string
	}
	tests := []struct {
		name      string
		fields    fields
		wantParts []string
	}{
		// git@gitforge.example/user/repo
		{"noProto1", fields{"", "git", "", "gitforge.example", 0, "hello/world.git"}, []string{"gitforge.example", "hello", "world"}},
		{"noProto2", fields{"", "git", "", "gitforge.example", 0, "hello/world"}, []string{"gitforge.example", "hello", "world"}},
		{"noProto3", fields{"", "git", "", "gitforge.example", 0, "hello/world/"}, []string{"gitforge.example", "hello", "world"}},
		{"noProto4", fields{"", "git", "", "gitforge.example", 0, "hello/world//"}, []string{"gitforge.example", "hello", "world"}},

		// ssh://git@gitforge.example/hello/world
		{"sshProto1", fields{"ssh", "git", "", "gitforge.example", 0, "hello/world.git"}, []string{"gitforge.example", "hello", "world"}},
		{"sshProto2", fields{"ssh", "git", "", "gitforge.example", 0, "hello/world"}, []string{"gitforge.example", "hello", "world"}},
		{"sshProto3", fields{"ssh", "git", "", "gitforge.example", 0, "hello/world/"}, []string{"gitforge.example", "hello", "world"}},
		{"sshProto4", fields{"ssh", "git", "", "gitforge.example", 0, "hello/world//"}, []string{"gitforge.example", "hello", "world"}},

		// user@server.com
		{"userServer1", fields{"", "user", "", "server.com", 0, "repository"}, []string{"server.com", "user", "repository"}},
		{"userServer2", fields{"", "user", "", "server.com", 0, "repository/"}, []string{"server.com", "user", "repository"}},
		{"userServer3", fields{"", "user", "", "server.com", 0, "repository//"}, []string{"server.com", "user", "repository"}},
		{"userServer4", fields{"", "user", "", "server.com", 0, "repository.git"}, []string{"server.com", "user", "repository"}},

		// ssh://user@server.com:1234
		{"userPort1", fields{"", "user", "", "server.com", 1234, "repository"}, []string{"server.com", "user", "repository"}},
		{"userPort2", fields{"", "user", "", "server.com", 1234, "repository/"}, []string{"server.com", "user", "repository"}},
		{"userPort3", fields{"", "user", "", "server.com", 1234, "repository//"}, []string{"server.com", "user", "repository"}},
		{"userPort4", fields{"", "user", "", "server.com", 1234, "repository.git"}, []string{"server.com", "user", "repository"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			url := env.URL{
				Scheme:   tt.fields.Scheme,
				User:     tt.fields.User,
				Password: tt.fields.Password,
				HostName: tt.fields.HostName,
				Port:     tt.fields.Port,
				Path:     tt.fields.Path,
			}
			if gotParts := url.Components(); !reflect.DeepEqual(gotParts, tt.wantParts) {
				t.Errorf("RepoURI.Components() = %v, want %v", gotParts, tt.wantParts)
			}
		})
	}
}

var benchComponentURLS = []env.URL{
	{"", "git", "", "gitforge.example", 0, "hello/world.git"},
	{"", "git", "", "gitforge.example", 0, "hello/world"},
	{"", "git", "", "gitforge.example", 0, "hello/world/"},
	{"", "git", "", "gitforge.example", 0, "hello/world//"},
	{"ssh", "git", "", "gitforge.example", 0, "hello/world.git"},
	{"ssh", "git", "", "gitforge.example", 0, "hello/world"},
	{"ssh", "git", "", "gitforge.example", 0, "hello/world/"},
	{"ssh", "git", "", "gitforge.example", 0, "hello/world//"},

	{"", "user", "", "server.com", 0, "repository"},
	{"", "user", "", "server.com", 0, "repository/"},
	{"", "user", "", "server.com", 0, "repository//"},
	{"", "user", "", "server.com", 0, "repository.git"},

	{"", "user", "", "server.com", 1234, "repository"},
	{"", "user", "", "server.com", 1234, "repository/"},
	{"", "user", "", "server.com", 1234, "repository//"},
	{"", "user", "", "server.com", 1234, "repository.git"},
}

func BenchmarkURL_Components(b *testing.B) {
	for b.Loop() {
		benchComponentURLS[0].Components()
		benchComponentURLS[1].Components()
		benchComponentURLS[2].Components()
		benchComponentURLS[3].Components()
		benchComponentURLS[4].Components()
		benchComponentURLS[5].Components()
		benchComponentURLS[6].Components()
		benchComponentURLS[7].Components()
		benchComponentURLS[8].Components()
		benchComponentURLS[9].Components()
		benchComponentURLS[10].Components()
		benchComponentURLS[11].Components()
		benchComponentURLS[12].Components()
		benchComponentURLS[13].Components()
		benchComponentURLS[14].Components()
		benchComponentURLS[15].Components()
	}
}

func TestComponentsOf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		s    string
		want []string
	}{
		{"ssh://host.xz/path/to/repo.git/", []string{"host.xz", "path", "to", "repo"}},
		{"ssh://user@host.xz/path/to/repo.git/", []string{"host.xz", "user", "path", "to", "repo"}},
		{"ssh://host.xz:1234/path/to/repo.git/", []string{"host.xz", "path", "to", "repo"}},
		{"ssh://user@host.xz:1234/path/to/repo.git/", []string{"host.xz", "user", "path", "to", "repo"}},
		{"git://host.xz/path/to/repo.git/", []string{"host.xz", "path", "to", "repo"}},
		{"git://host.xz:1234/path/to/repo.git/", []string{"host.xz", "path", "to", "repo"}},
		{"host.xz:path/to/repo.git/", []string{"host.xz", "path", "to", "repo"}},
		{"user@host.xz:path/to/repo.git/", []string{"host.xz", "user", "path", "to", "repo"}},
		{"user@host.xz:path/to/repo", []string{"host.xz", "user", "path", "to", "repo"}},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			t.Parallel()

			if got := env.ComponentsOf(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComponentsOf() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func BenchmarkComponentsOf(b *testing.B) {
	for b.Loop() {
		env.ComponentsOf("ssh://host.xz/path/to/repo.git/")
		env.ComponentsOf("ssh://user@host.xz/path/to/repo.git/")
		env.ComponentsOf("ssh://host.xz:1234/path/to/repo.git/")
		env.ComponentsOf("ssh://user@host.xz:1234/path/to/repo.git/")
		env.ComponentsOf("git://host.xz/path/to/repo.git/")
		env.ComponentsOf("git://host.xz:1234/path/to/repo.git/")
		env.ComponentsOf("host.xz:path/to/repo.git/")
		env.ComponentsOf("user@host.xz:path/to/repo.git/")
		env.ComponentsOf("user@host.xz:path/to/repo")
	}
}

func TestRepoURL_Canonical(t *testing.T) {
	t.Parallel()

	type fields struct {
		Scheme   string
		User     string
		Password string
		HostName string
		Port     uint16
		Path     string
	}
	type args struct {
		cSpec string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantCanonical string
	}{
		{"Treat one component special", fields{"", "user", "", "server.com", 1234, "repository"}, args{"git@^:$.git"}, "git@server.com:user/repository.git"},
		{"Treat two components special", fields{"", "user", "", "server.com", 1234, "repository"}, args{"ssh://%@^/$.git"}, "ssh://user@server.com/repository.git"},
		{"Empty specification string", fields{Scheme: "", User: "user", Password: "", HostName: "server.com", Port: 1234, Path: "repository"}, args{""}, "server.com/user/repository"},
		{"Return original url", fields{"", "user", "", "server.com", 1234, "repository"}, args{"$$"}, "user@server.com:1234:repository"},
		{"Drop a component", fields{"", "user", "", "server.com", 1234, "repository"}, args{"git@!other_server.com:$.git"}, "git@other_server.com:user/repository.git"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rURL := env.URL{
				Scheme:   tt.fields.Scheme,
				User:     tt.fields.User,
				Password: tt.fields.Password,
				HostName: tt.fields.HostName,
				Port:     tt.fields.Port,
				Path:     tt.fields.Path,
			}
			if gotCanonical := rURL.Canonical(tt.args.cSpec); gotCanonical != tt.wantCanonical {
				t.Errorf("RepoURL.Canonical() = %v, want %v", gotCanonical, tt.wantCanonical)
			}
		})
	}
}
