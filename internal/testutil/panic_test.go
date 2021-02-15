package testutil

import (
	"reflect"
	"testing"
)

func TestDoesPanic(t *testing.T) {
	type args struct {
		f func()
	}
	tests := []struct {
		name          string
		args          args
		wantPaniced   bool
		wantRecovered interface{}
	}{
		{"not panicing f", args{func() {}}, false, nil},
		{"panicing with message", args{func() { panic("message") }}, true, "message"},
		{"panicing with nil", args{func() { panic(nil) }}, true, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPaniced, gotRecovered := DoesPanic(tt.args.f)
			if gotPaniced != tt.wantPaniced {
				t.Errorf("DoesPanic() gotPaniced = %v, want %v", gotPaniced, tt.wantPaniced)
			}
			if !reflect.DeepEqual(gotRecovered, tt.wantRecovered) {
				t.Errorf("DoesPanic() gotRecovered = %v, want %v", gotRecovered, tt.wantRecovered)
			}
		})
	}
}
