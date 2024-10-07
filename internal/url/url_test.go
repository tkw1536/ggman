package url

//spellchecker:words math strconv testing
import (
	"math"
	"strconv"
	"testing"
)

func Test_ParsePort(t *testing.T) {
	type args struct {
		portString string
	}
	tests := []struct {
		name     string
		args     args
		wantPort uint16
		wantErr  bool
	}{
		{"parsing zero port", args{"0"}, 0, false},
		{"parsing empty port", args{""}, 0, true},
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
			gotPort, err := ParsePort(tt.args.portString)
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
	for port := 0; port <= maxValid; port++ {
		gotPort, err := ParsePort(strconv.Itoa(port))
		if gotPort != uint16(port) {
			t.Errorf("ParsePort(%d) got port = %d", port, gotPort)
		}
		if err != nil {
			t.Errorf("ParsePort(%d) got error = %v, want error = nil", port, err)
		}
	}
	for port := maxValid + 1; port <= maxPortTest; port++ {
		gotPort, err := ParsePort(strconv.Itoa(port))
		if gotPort != 0 {
			t.Errorf("ParsePort(%d) got port = %d", port, gotPort)
		}
		if err == nil {
			t.Errorf("ParsePort(%d) got error = nil, but want not nil", port)
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

func Test_SplitURLScheme(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name       string
		args       args
		wantScheme string
		wantRest   string
	}{
		{name: "http scheme is valid", args: args{"http://rest"}, wantScheme: "http", wantRest: "rest"},
		{name: "https scheme is valid", args: args{"https://rest"}, wantScheme: "https", wantRest: "rest"},
		{name: "git scheme is valid", args: args{"git://rest"}, wantScheme: "git", wantRest: "rest"},
		{name: "ssh scheme is valid", args: args{"ssh://rest"}, wantScheme: "ssh", wantRest: "rest"},
		{name: "ssh2 scheme is valid", args: args{"ssh2://rest"}, wantScheme: "ssh2", wantRest: "rest"},
		{name: "file scheme is valid", args: args{"file://rest"}, wantScheme: "file", wantRest: "rest"},
		{name: "combined with + scheme is valid", args: args{"ssh+git://rest"}, wantScheme: "ssh+git", wantRest: "rest"},
		{name: "combined with - scheme is valid", args: args{"ssh-git://rest"}, wantScheme: "ssh-git", wantRest: "rest"},
		{name: "combined with . scheme is valid", args: args{"ssh.git://rest"}, wantScheme: "ssh.git", wantRest: "rest"},

		{name: "valid scheme without valid ending", args: args{"valid"}, wantScheme: "", wantRest: "valid"},
		{name: "valid scheme with only :", args: args{"valid:"}, wantScheme: "", wantRest: "valid:"},
		{name: "valid scheme with only one /", args: args{"valid:/"}, wantScheme: "", wantRest: "valid:/"},
		{name: "valid scheme without rest", args: args{"valid://"}, wantScheme: "valid", wantRest: ""},

		{name: "empty input passed through", args: args{""}, wantScheme: "", wantRest: ""},
		{name: "numerical scheme is invalid", args: args{"01234"}, wantScheme: "", wantRest: "01234"},
		{name: "numerical starting scheme is invalid", args: args{"01git"}, wantScheme: "", wantRest: "01git"},
		{name: "non-scheme is invalid", args: args{"://"}, wantScheme: "", wantRest: "://"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScheme, gotRest := SplitURLScheme(tt.args.input)
			if gotScheme != tt.wantScheme {
				t.Errorf("SplitURLScheme() scheme = %v, want %v", gotScheme, tt.wantScheme)
			}
			if gotRest != tt.wantRest {
				t.Errorf("SplitURLScheme() rest = %v, want %v", gotRest, tt.wantRest)
			}
		})
	}
}
