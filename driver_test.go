package dialsrv

import (
	"context"
	"net"
	"strings"
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
	if got := err.Error(); !strings.HasPrefix(got, "lookup ") || !strings.HasSuffix(got, ": no such host"){
		t.Errorf("unexpected error:\nwant=%s\n got=%s", "lookup .*: no such host", got)
	}
}
