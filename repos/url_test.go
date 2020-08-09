package repos

import (
	"reflect"
	"testing"
)

func TestParseRepoURL(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name     string
		args     args
		wantRepo *RepoURL
	}{
		// ssh://[user@]host.xz[:port]/path/to/repo.git/
		{
			"ssh",
			args{"ssh://host.xz/path/to/repo.git/"},
			&RepoURL{"ssh", "", "", "host.xz", 0, "path/to/repo.git/"},
		},
		{
			"sshUser",
			args{"ssh://user@host.xz/path/to/repo.git/"},
			&RepoURL{"ssh", "user", "", "host.xz", 0, "path/to/repo.git/"},
		},
		{
			"sshPort",
			args{"ssh://host.xz:1234/path/to/repo.git/"},
			&RepoURL{"ssh", "", "", "host.xz", 1234, "path/to/repo.git/"},
		},
		{
			"sshUserPort",
			args{"ssh://user@host.xz:1234/path/to/repo.git/"},
			&RepoURL{"ssh", "user", "", "host.xz", 1234, "path/to/repo.git/"},
		},

		// git://host.xz[:port]/path/to/repo.git/
		{
			"git",
			args{"git://host.xz/path/to/repo.git/"},
			&RepoURL{"git", "", "", "host.xz", 0, "path/to/repo.git/"},
		},

		{
			"gitPort",
			args{"git://host.xz:1234/path/to/repo.git/"},
			&RepoURL{"git", "", "", "host.xz", 1234, "path/to/repo.git/"},
		},

		//  [user@]host.xz:path/to/repo.git/
		{
			"noProto",
			args{"host.xz:path/to/repo.git/"},
			&RepoURL{"", "", "", "host.xz", 0, "path/to/repo.git/"},
		},
		{
			"noProtoUser",
			args{"user@host.xz:path/to/repo.git/"},
			&RepoURL{"", "user", "", "host.xz", 0, "path/to/repo.git/"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepo := ParseRepoURL(tt.args.s)
			if !reflect.DeepEqual(gotRepo, tt.wantRepo) {
				t.Errorf("ParseRepoURL() = %v, want %v", gotRepo, tt.wantRepo)
			}
		})
	}
}
