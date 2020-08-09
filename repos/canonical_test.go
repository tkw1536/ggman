package repos

import "testing"

func TestRepoURI_Canonical(t *testing.T) {
	type fields struct {
		Scheme   string
		User     string
		Password string
		HostName string
		Port     int
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
			rURI := &RepoURL{
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
