package dialsrv

import (
	"reflect"
	"testing"
)

func TestParseAddr(t *testing.T) {
	for _, d := range []struct {
		n, a string
		fa   FlavoredAddr
	}{
		{"tcp", "srv+myservice+example.com",
			FlavoredAddr{"tcp", "myservice", "tcp", "example.com"}},
		{"udp", "srv+myservice+example.com",
			FlavoredAddr{"udp", "myservice", "udp", "example.com"}},
		{"tcp", "srv+myapi+example.com",
			FlavoredAddr{"tcp", "myapi", "tcp", "example.com"}},
		{"tcp", "srv+myservice+foo.example.org",
			FlavoredAddr{"tcp", "myservice", "tcp", "foo.example.org"}},
		{"tcp", "srv+example.com",
			FlavoredAddr{"tcp", "", "", "example.com"}},
	} {
		act := parseAddr(d.n, d.a)
		if !reflect.DeepEqual(act, &d.fa) {
			t.Errorf("unexpected parse %s, %s: %#v", d.n, d.a, act)
		}
	}
}

func TestParseAddrNil(t *testing.T) {
	for _, d := range []struct {
		n, a string
	}{
		{"tcp", "example.com"},
		{"udp", "example.com"},
		{"tcp", "foo.example.org"},
		{"tcp", "foo.example.com"},
	} {
		act := parseAddr(d.n, d.a)
		if act != nil {
			t.Errorf("unexpected non-nil %s, %s: %#v", d.n, d.a, act)
		}
	}
}
