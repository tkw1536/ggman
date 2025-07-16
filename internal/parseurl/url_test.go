//spellchecker:words parseurl
package parseurl_test

//spellchecker:words math strconv testing ggman internal parseurl
import (
	"math"
	"strconv"
	"testing"

	"go.tkw01536.de/ggman/internal/parseurl"
)

func Test_ParsePort(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			gotPort, err := parseurl.ParsePort(tt.args.portString)
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

// constants copied over from the url package proper
// so that we don't need to make them public!
const (
	maxValid   = math.MaxUint16  // maximal port (as a number)
	maxPortStr = "65535"         // maximal port (as a string)
	maxPortLen = len(maxPortStr) // maximal port length
)

// maximal port number to test
// this should really be a const, but there's no pow in that.
var maxPortTest = int(math.Pow(10, float64(maxPortLen)) - 1)

func Test_ParsePort_all(t *testing.T) {
	t.Parallel()

	for port := 0; port <= maxValid; port++ {
		t.Run(strconv.Itoa(port), func(t *testing.T) {
			t.Parallel()

			if port < 0 || port > math.MaxUint16 {
				// bounds check to make linter happy!
				panic("never reached")
			}

			gotPort, err := parseurl.ParsePort(strconv.Itoa(port))
			if gotPort != uint16(port) {
				t.Errorf("ParsePort(%d) got port = %d", port, gotPort)
			}
			if err != nil {
				t.Errorf("ParsePort(%d) got error = %v, want error = nil", port, err)
			}
		})
	}

	for port := maxValid + 1; port <= maxPortTest; port++ {
		t.Run(strconv.Itoa(port), func(t *testing.T) {
			t.Parallel()

			gotPort, err := parseurl.ParsePort(strconv.Itoa(port))
			if gotPort != 0 {
				t.Errorf("ParsePort(%d) got port = %d", port, gotPort)
			}
			if err == nil {
				t.Errorf("ParsePort(%d) got error = nil, but want not nil", port)
			}
		})
	}
}

func Benchmark_ParsePort(b *testing.B) {
	for b.Loop() {
		// ignore all the errors, cause we're benchmarking!
		_, _ = parseurl.ParsePort("0")
		_, _ = parseurl.ParsePort("80")
		_, _ = parseurl.ParsePort("65535")
		_, _ = parseurl.ParsePort(" 8080")
		_, _ = parseurl.ParsePort("-1")
		_, _ = parseurl.ParsePort("65536")
		_, _ = parseurl.ParsePort("+123")
		_, _ = parseurl.ParsePort("aaaa")
	}
}

func Test_SplitScheme(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			gotScheme, gotRest := parseurl.SplitScheme(tt.args.input)
			if gotScheme != tt.wantScheme {
				t.Errorf("SplitScheme() scheme = %v, want %v", gotScheme, tt.wantScheme)
			}
			if gotRest != tt.wantRest {
				t.Errorf("SplitScheme() rest = %v, want %v", gotRest, tt.wantRest)
			}
		})
	}
}

func Benchmark_SplitScheme(b *testing.B) {
	for b.Loop() {
		_, _ = parseurl.SplitScheme("http://rest")
		_, _ = parseurl.SplitScheme("https://rest")
		_, _ = parseurl.SplitScheme("git://rest")
		_, _ = parseurl.SplitScheme("ssh://rest")
		_, _ = parseurl.SplitScheme("ssh2://rest")
		_, _ = parseurl.SplitScheme("file://rest")
		_, _ = parseurl.SplitScheme("ssh+git://rest")
		_, _ = parseurl.SplitScheme("ssh-git://rest")
		_, _ = parseurl.SplitScheme("ssh.git://rest")
		_, _ = parseurl.SplitScheme("valid")
		_, _ = parseurl.SplitScheme("valid:")
		_, _ = parseurl.SplitScheme("valid:/")
		_, _ = parseurl.SplitScheme("valid://")
		_, _ = parseurl.SplitScheme("")
		_, _ = parseurl.SplitScheme("01234")
		_, _ = parseurl.SplitScheme("01git")
		_, _ = parseurl.SplitScheme("://")
	}
}
