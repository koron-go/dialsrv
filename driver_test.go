package dialsrv

import (
	"context"
	"net"
	"testing"
)

func TestDialerDriver(t *testing.T) {
	driver := &netDialerDriver{&net.Dialer{}}
	cname, addrs, err := driver.LookupSRV(context.Background(), "subservice", "tcp", "example.com")
	if cname != "" {
		t.Errorf("unexpected cname=%s", cname)
	}
	if len(addrs) != 0 {
		t.Errorf("unexpected addrs=%+v", addrs)
	}
	if want, got := "lookup example.com: no such host", err.Error(); got != want {
		t.Errorf("unexpected error:\nwant=%s\n got=%s", want, got)
	}
}
