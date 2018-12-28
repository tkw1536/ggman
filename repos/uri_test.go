package repos

import (
	"reflect"
	"testing"
)

func TestNewRepoURI(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name     string
		args     args
		wantRepo *RepoURI
		wantErr  bool
	}{
		// ssh://[user@]host.xz[:port]/path/to/repo.git/
		{
			"ssh",
			args{"ssh://host.xz/path/to/repo.git/"},
			&RepoURI{"ssh", "", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},
		{
			"sshUser",
			args{"ssh://user@host.xz/path/to/repo.git/"},
			&RepoURI{"ssh", "user", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},
		{
			"sshPort",
			args{"ssh://host.xz:1234/path/to/repo.git/"},
			&RepoURI{"ssh", "", "", "host.xz", 1234, "path/to/repo.git/"},
			false,
		},
		{
			"sshUserPort",
			args{"ssh://user@host.xz:1234/path/to/repo.git/"},
			&RepoURI{"ssh", "user", "", "host.xz", 1234, "path/to/repo.git/"},
			false,
		},

		// git://host.xz[:port]/path/to/repo.git/
		{
			"git",
			args{"git://host.xz/path/to/repo.git/"},
			&RepoURI{"git", "", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},

		{
			"gitPort",
			args{"git://host.xz:1234/path/to/repo.git/"},
			&RepoURI{"git", "", "", "host.xz", 1234, "path/to/repo.git/"},
			false,
		},

		//  [user@]host.xz:path/to/repo.git/
		{
			"noProto",
			args{"host.xz:path/to/repo.git/"},
			&RepoURI{"", "", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},
		{
			"noProtoUser",
			args{"user@host.xz:path/to/repo.git/"},
			&RepoURI{"", "user", "", "host.xz", 0, "path/to/repo.git/"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepo, err := NewRepoURI(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRepoURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRepo, tt.wantRepo) {
				t.Errorf("NewRepoURI() = %v, want %v", gotRepo, tt.wantRepo)
			}
		})
	}
}

func TestRepoURI_Components(t *testing.T) {
	type fields struct {
		Scheme   string
		User     string
		Password string
		HostName string
		Port     int
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
			rURI := &RepoURI{
				Scheme:   tt.fields.Scheme,
				User:     tt.fields.User,
				Password: tt.fields.Password,
				HostName: tt.fields.HostName,
				Port:     tt.fields.Port,
				Path:     tt.fields.Path,
			}
			if gotParts := rURI.Components(); !reflect.DeepEqual(gotParts, tt.wantParts) {
				t.Errorf("RepoURI.Components() = %v, want %v", gotParts, tt.wantParts)
			}
		})
	}
}
