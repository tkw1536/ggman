package env

import (
	"reflect"
	"testing"
)

func TestParseURL(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name     string
		args     args
		wantRepo URL
	}{
		// ssh://[user@]host.xz[:port]/path/to/repo.git/
		{
			"ssh",
			args{"ssh://host.xz/path/to/repo.git/"},
			URL{"ssh", "", "", "host.xz", 0, "path/to/repo.git/"},
		},
		{
			"sshUser",
			args{"ssh://user@host.xz/path/to/repo.git/"},
			URL{"ssh", "user", "", "host.xz", 0, "path/to/repo.git/"},
		},
		{
			"sshPort",
			args{"ssh://host.xz:1234/path/to/repo.git/"},
			URL{"ssh", "", "", "host.xz", 1234, "path/to/repo.git/"},
		},
		{
			"sshUserPort",
			args{"ssh://user@host.xz:1234/path/to/repo.git/"},
			URL{"ssh", "user", "", "host.xz", 1234, "path/to/repo.git/"},
		},

		// git://host.xz[:port]/path/to/repo.git/
		{
			"git",
			args{"git://host.xz/path/to/repo.git/"},
			URL{"git", "", "", "host.xz", 0, "path/to/repo.git/"},
		},

		{
			"gitPort",
			args{"git://host.xz:1234/path/to/repo.git/"},
			URL{"git", "", "", "host.xz", 1234, "path/to/repo.git/"},
		},

		//  [user@]host.xz:path/to/repo.git/
		{
			"noProto",
			args{"host.xz:path/to/repo.git/"},
			URL{"", "", "", "host.xz", 0, "path/to/repo.git/"},
		},
		{
			"noProtoUser",
			args{"user@host.xz:path/to/repo.git/"},
			URL{"", "user", "", "host.xz", 0, "path/to/repo.git/"},
		},

		// local paths
		{
			"localFile",
			args{"file:///path/to/somewhere"},
			URL{"file", "", "", "", 0, "path/to/somewhere"},
		},

		{
			"localPath",
			args{"/path/to/somewhere"},
			URL{"", "", "", "", 0, "path/to/somewhere"},
		},

		{
			"localRelPath",
			args{"../some/relative/path"},
			URL{"", "", "", "..", 0, "some/relative/path"},
		},

		{
			"localRelPath2",
			args{"./some/relative/path"},
			URL{"", "", "", ".", 0, "some/relative/path"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepo := ParseURL(tt.args.s)
			if !reflect.DeepEqual(gotRepo, tt.wantRepo) {
				t.Errorf("ParseRepoURL() = %v, want %v", gotRepo, tt.wantRepo)
			}
		})
	}
}

func Benchmark_ParseRepoURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseURL("ssh://host.xz/path/to/repo.git/")
		ParseURL("ssh://user@host.xz/path/to/repo.git/")
		ParseURL("ssh://host.xz:1234/path/to/repo.git/")
		ParseURL("ssh://user@host.xz:1234/path/to/repo.git/")
		ParseURL("git://host.xz/path/to/repo.git/")
		ParseURL("git://host.xz:1234/path/to/repo.git/")
		ParseURL("host.xz:path/to/repo.git/")
		ParseURL("user@host.xz:path/to/repo.git/")
	}
}

func TestURL_IsLocal(t *testing.T) {
	tests := []struct {
		name string
		url  URL
		want bool
	}{
		{
			"ssh",
			URL{"ssh", "", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},
		{
			"sshUser",
			URL{"ssh", "user", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},
		{
			"sshPort",
			URL{"ssh", "", "", "host.xz", 1234, "path/to/repo.git/"},
			false,
		},
		{
			"sshUserPort",
			URL{"ssh", "user", "", "host.xz", 1234, "path/to/repo.git/"},
			false,
		},

		// git://host.xz[:port]/path/to/repo.git/
		{
			"git",
			URL{"git", "", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},

		{
			"gitPort",
			URL{"git", "", "", "host.xz", 1234, "path/to/repo.git/"},
			false,
		},

		//  [user@]host.xz:path/to/repo.git/
		{
			"noProto",
			URL{"", "", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},
		{
			"noProtoUser",
			URL{"", "user", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},

		// local paths
		{
			"localFile",
			URL{"file", "", "", "", 0, "path/to/somewhere"},
			true,
		},

		{
			"localPath",
			URL{"", "", "", "", 0, "path/to/somewhere"},
			true,
		},

		{
			"localRelPath",
			URL{"", "", "", "..", 0, "some/relative/path"},
			true,
		},

		{
			"localRelPath2",
			URL{"", "", "", ".", 0, "some/relative/path"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.url.IsLocal(); got != tt.want {
				t.Errorf("URL.IsLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURL_Components(t *testing.T) {
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
		{"noproto1", fields{"", "git", "", "github.com", 0, "hello/world.git"}, []string{"github.com", "hello", "world"}},
		{"noproto2", fields{"", "git", "", "github.com", 0, "hello/world"}, []string{"github.com", "hello", "world"}},
		{"noproto3", fields{"", "git", "", "github.com", 0, "hello/world/"}, []string{"github.com", "hello", "world"}},
		{"noproto4", fields{"", "git", "", "github.com", 0, "hello/world//"}, []string{"github.com", "hello", "world"}},

		// ssh://git@github.com/hello/world
		{"sshproto1", fields{"ssh", "git", "", "github.com", 0, "hello/world.git"}, []string{"github.com", "hello", "world"}},
		{"sshproto2", fields{"ssh", "git", "", "github.com", 0, "hello/world"}, []string{"github.com", "hello", "world"}},
		{"sshproto3", fields{"ssh", "git", "", "github.com", 0, "hello/world/"}, []string{"github.com", "hello", "world"}},
		{"sshproto4", fields{"ssh", "git", "", "github.com", 0, "hello/world//"}, []string{"github.com", "hello", "world"}},

		// user@server.com
		{"userserver1", fields{"", "user", "", "server.com", 0, "repository"}, []string{"server.com", "user", "repository"}},
		{"userserver2", fields{"", "user", "", "server.com", 0, "repository/"}, []string{"server.com", "user", "repository"}},
		{"userserver3", fields{"", "user", "", "server.com", 0, "repository//"}, []string{"server.com", "user", "repository"}},
		{"userserver4", fields{"", "user", "", "server.com", 0, "repository.git"}, []string{"server.com", "user", "repository"}},

		// ssh://user@server.com:1234
		{"userport1", fields{"", "user", "", "server.com", 1234, "repository"}, []string{"server.com", "user", "repository"}},
		{"userport2", fields{"", "user", "", "server.com", 1234, "repository/"}, []string{"server.com", "user", "repository"}},
		{"userport3", fields{"", "user", "", "server.com", 1234, "repository//"}, []string{"server.com", "user", "repository"}},
		{"userport4", fields{"", "user", "", "server.com", 1234, "repository.git"}, []string{"server.com", "user", "repository"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := URL{
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

var benchComponentURLS = []URL{
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
	for i := 0; i < b.N; i++ {
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
			if got := ComponentsOf(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComponentsOf() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func BenchmarkComponentsOf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ComponentsOf("ssh://host.xz/path/to/repo.git/")
		ComponentsOf("ssh://user@host.xz/path/to/repo.git/")
		ComponentsOf("ssh://host.xz:1234/path/to/repo.git/")
		ComponentsOf("ssh://user@host.xz:1234/path/to/repo.git/")
		ComponentsOf("git://host.xz/path/to/repo.git/")
		ComponentsOf("git://host.xz:1234/path/to/repo.git/")
		ComponentsOf("host.xz:path/to/repo.git/")
		ComponentsOf("user@host.xz:path/to/repo.git/")
		ComponentsOf("user@host.xz:path/to/repo")
	}
}

func TestRepoURI_Canonical(t *testing.T) {
	type fields struct {
		Scheme   string
		User     string
		Password string
		HostName string
		Port     uint16
		Path     string
	}
	type args struct {
		cspec string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantCanonical string
	}{
		{"Treat one component special", fields{"", "user", "", "server.com", 1234, "repository"}, args{"git@^:$.git"}, "git@server.com:user/repository.git"},
		{"Treat two components special", fields{"", "user", "", "server.com", 1234, "repository"}, args{"ssh://%@^/$.git"}, "ssh://user@server.com/repository.git"},
		{"Empty specifcation string", fields{"", "user", "", "server.com", 1234, "repository"}, args{""}, "server.com/user/repository"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rURI := &URL{
				Scheme:   tt.fields.Scheme,
				User:     tt.fields.User,
				Password: tt.fields.Password,
				HostName: tt.fields.HostName,
				Port:     tt.fields.Port,
				Path:     tt.fields.Path,
			}
			if gotCanonical := rURI.Canonical(tt.args.cspec); gotCanonical != tt.wantCanonical {
				t.Errorf("RepoURI.Canonical() = %v, want %v", gotCanonical, tt.wantCanonical)
			}
		})
	}
}
