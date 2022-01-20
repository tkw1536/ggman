package url

import (
	"math"
	"strconv"
	"testing"
)

func Test_ParsePort(t *testing.T) {
	type args struct {
		portstring string
	}
	tests := []struct {
		name     string
		args     args
		wantPort uint16
		wantErr  bool
	}{
		{"parsing zero port", args{"0"}, 0, false},
		{"parsing valid port", args{"80"}, 80, false},
		{"parsing high port", args{"65535"}, 65535, false},

		{"parsing port with space", args{" 8080"}, 0, true},
		{"parsing negative port", args{"-1"}, 0, true},
		{"parsing too high port", args{"65536"}, 0, true},
		{"parsing port with '+' in it", args{"+123"}, 0, true},
		{"parsing non-numeric port", args{"aaaa"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPort, err := ParsePort(tt.args.portstring)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPort != tt.wantPort {
				t.Errorf("parsePort() = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}

var maxPortTest = int(math.Pow(10, float64(maxPortLen)) - 1)

func Test_ParsePort_all(t *testing.T) {
	for port := 0; port <= maxValidPort; port++ {
		gotPort, err := ParsePort(strconv.Itoa(port))
		if gotPort != uint16(port) {
			t.Errorf("ParsePort(%d) got port = %d", port, gotPort)
		}
		if err != nil {
			t.Errorf("ParsePort(%d) got error = %v, want error = nil", port, err)
		}
	}
	for port := maxValidPort + 1; port <= maxPortTest; port++ {
		gotPort, err := ParsePort(strconv.Itoa(port))
		if gotPort != 0 {
			t.Errorf("ParsePort(%d) got port = %d", port, gotPort)
		}
		if err != errInvalidRange {
			t.Errorf("ParsePort(%d) got error = %v, want error = errInvalidRange", port, err)
		}
	}
}

func Benchmark_ParsePort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParsePort("0")
		ParsePort("80")
		ParsePort("65535")
		ParsePort(" 8080")
		ParsePort("-1")
		ParsePort("65536")
		ParsePort("+123")
		ParsePort("aaaa")
	}
}

func Test_IsValidURLScheme(t *testing.T) {
	type args struct {
		scheme string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"http scheme is valid", args{"http"}, true},
		{"https scheme is valid", args{"https"}, true},
		{"git scheme is valid", args{"git"}, true},
		{"ssh scheme is valid", args{"ssh"}, true},
		{"file scheme is valid", args{"file"}, true},
		{"combined with + scheme is valid", args{"ssh+git"}, true},
		{"combined with - scheme is valid", args{"ssh-git"}, true},
		{"combined with . scheme is valid", args{"ssh.git"}, true},

		{"empty scheme is invalid", args{""}, false},
		{"numerical scheme is invalid", args{"01234"}, false},
		{"numerical starting scheme is invalid", args{"01git"}, false},
		{"non-scheme is invalid", args{"://"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidURLScheme(tt.args.scheme); got != tt.want {
				t.Errorf("validateScheme() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_IsValidURLScheme(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsValidURLScheme("0http")
		IsValidURLScheme("http")
	}
}
