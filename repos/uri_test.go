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
