package env_test

//spellchecker:words reflect testing github ggman
import (
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/env"
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

func Benchmark_ParseRepoURL(b *testing.B) {
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

func TestURL_IsLocal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url  env.URL
		want bool
	}{
		{
			"ssh",
			env.URL{"ssh", "", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},
		{
			"sshUser",
			env.URL{"ssh", "user", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},
		{
			"sshPort",
			env.URL{"ssh", "", "", "host.xz", 1234, "path/to/repo.git/"},
			false,
		},
		{
			"sshUserPort",
			env.URL{"ssh", "user", "", "host.xz", 1234, "path/to/repo.git/"},
			false,
		},

		// git://host.xz[:port]/path/to/repo.git/
		{
			"git",
			env.URL{"git", "", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},

		{
			"gitPort",
			env.URL{"git", "", "", "host.xz", 1234, "path/to/repo.git/"},
			false,
		},

		//  [user@]host.xz:path/to/repo.git/
		{
			"noProto",
			env.URL{"", "", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},
		{
			"noProtoUser",
			env.URL{"", "user", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},

		// local paths
		{
			"localFile",
			env.URL{"file", "", "", "", 0, "path/to/somewhere"},
			true,
		},

		{
			"localPath",
			env.URL{"", "", "", "", 0, "path/to/somewhere"},
			true,
		},

		{
			"localRelPath",
			env.URL{"", "", "", "..", 0, "some/relative/path"},
			true,
		},

		{
			"localRelPath2",
			env.URL{"", "", "", ".", 0, "some/relative/path"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.url.IsLocal(); got != tt.want {
				t.Errorf("URL.IsLocal() = %v, want %v", got, tt.want)
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
		// git@github.com/user/repo
		{"noProto1", fields{"", "git", "", "github.com", 0, "hello/world.git"}, []string{"github.com", "hello", "world"}},
		{"noProto2", fields{"", "git", "", "github.com", 0, "hello/world"}, []string{"github.com", "hello", "world"}},
		{"noProto3", fields{"", "git", "", "github.com", 0, "hello/world/"}, []string{"github.com", "hello", "world"}},
		{"noProto4", fields{"", "git", "", "github.com", 0, "hello/world//"}, []string{"github.com", "hello", "world"}},

		// ssh://git@github.com/hello/world
		{"sshProto1", fields{"ssh", "git", "", "github.com", 0, "hello/world.git"}, []string{"github.com", "hello", "world"}},
		{"sshProto2", fields{"ssh", "git", "", "github.com", 0, "hello/world"}, []string{"github.com", "hello", "world"}},
		{"sshProto3", fields{"ssh", "git", "", "github.com", 0, "hello/world/"}, []string{"github.com", "hello", "world"}},
		{"sshProto4", fields{"ssh", "git", "", "github.com", 0, "hello/world//"}, []string{"github.com", "hello", "world"}},

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
	{"", "git", "", "github.com", 0, "hello/world.git"},
	{"", "git", "", "github.com", 0, "hello/world"},
	{"", "git", "", "github.com", 0, "hello/world/"},
	{"", "git", "", "github.com", 0, "hello/world//"},
	{"ssh", "git", "", "github.com", 0, "hello/world.git"},
	{"ssh", "git", "", "github.com", 0, "hello/world"},
	{"ssh", "git", "", "github.com", 0, "hello/world/"},
	{"ssh", "git", "", "github.com", 0, "hello/world//"},

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
