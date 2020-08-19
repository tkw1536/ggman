package repos

import (
	"reflect"
	"testing"
)

func TestRepoURI_Components(t *testing.T) {
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
			rURI := &URL{
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
