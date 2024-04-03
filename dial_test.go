package dialsrv

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestParseAddr(t *testing.T) {
	for _, d := range []struct {
		n, a string
		fa   FlavoredAddr
		str  string
	}{
		{
			"tcp", "srv+myservice+example.com",
			FlavoredAddr{"tcp", "myservice", "tcp", "example.com"},
			"_myservice._tcp.example.com",
		},
		{
			"udp", "srv+myservice+example.com",
			FlavoredAddr{"udp", "myservice", "udp", "example.com"},
			"_myservice._udp.example.com",
		},
		{
			"tcp", "srv+myapi+example.com",
			FlavoredAddr{"tcp", "myapi", "tcp", "example.com"},
			"_myapi._tcp.example.com",
		},
		{
			"tcp", "srv+myservice+foo.example.org",
			FlavoredAddr{"tcp", "myservice", "tcp", "foo.example.org"},
			"_myservice._tcp.foo.example.org",
		},
		{
			"tcp", "srv+example.com",
			FlavoredAddr{"tcp", "", "", "example.com"},
			"example.com",
		},
	} {
		act := parseAddr(d.n, d.a)
		if !reflect.DeepEqual(act, &d.fa) {
			t.Errorf("unexpected parse %s, %s: %#v", d.n, d.a, act)
		}
		if want, got := d.str, act.String(); want != got {
			t.Errorf("unexpected string:\nwant=%s\n got=%s", want, got)
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

type testConn struct {
	network string
	address string
}

func (*testConn) Read([]byte) (int, error)         { return 0, nil }
func (*testConn) Write([]byte) (int, error)        { return 0, nil }
func (*testConn) Close() error                     { return nil }
func (*testConn) LocalAddr() net.Addr              { return nil }
func (*testConn) RemoteAddr() net.Addr             { return nil }
func (*testConn) SetDeadline(time.Time) error      { return nil }
func (*testConn) SetReadDeadline(time.Time) error  { return nil }
func (*testConn) SetWriteDeadline(time.Time) error { return nil }

type dialContextParams struct {
	network string
	address string
}

type dialContextResults struct {
	conn net.Conn
	err  error
}

type lookupSRVParams struct {
	service string
	proto   string
	name    string
}

type lookupSRVResults struct {
	cname string
	addrs []*net.SRV
	err   error
}

type testDriver struct {
	dialContextParams  *dialContextParams
	dialContextResults *dialContextResults
	lookupSRVParams    *lookupSRVParams
	lookupSRVResults   *lookupSRVResults
}

func (d *testDriver) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	d.dialContextParams = &dialContextParams{
		network: network,
		address: address,
	}
	if d.dialContextResults == nil {
		return &testConn{
			network: network,
			address: address,
		}, nil
	}
	return d.dialContextResults.conn, d.dialContextResults.err
}

func (d *testDriver) LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
	d.lookupSRVParams = &lookupSRVParams{
		service: service,
		proto:   proto,
		name:    name,
	}
	if d.lookupSRVResults == nil {
		target := name
		if service != "" {
			target = service + "." + name
		}
		return "sample", []*net.SRV{
			{Target: target, Port: 1234, Priority: 1, Weight: 100},
		}, nil
	}
	return d.lookupSRVResults.cname, d.lookupSRVResults.addrs, d.lookupSRVResults.err
}

var _ driver = (*testDriver)(nil)

func TestDial(t *testing.T) {
	for _, c := range []struct {
		network          string
		address          string
		want             dialContextResults
		wantDialParams   *dialContextParams
		wantLookupParams *lookupSRVParams
	}{
		{ // without "srv+" prefix
			"tcp",
			"example.com",
			dialContextResults{&testConn{"tcp", "example.com"}, nil},
			&dialContextParams{"tcp", "example.com"},
			nil,
		},
		{ // with simple "srv+"
			"tcp",
			"srv+example.com",
			dialContextResults{&testConn{"tcp", "example.com:1234"}, nil},
			&dialContextParams{"tcp", "example.com:1234"},
			&lookupSRVParams{"", "", "example.com"},
		},
		{ // with "srv+" and "ldap" service
			"tcp",
			"srv+ldap+example.com",
			dialContextResults{&testConn{"tcp", "ldap.example.com:1234"}, nil},
			&dialContextParams{"tcp", "ldap.example.com:1234"},
			&lookupSRVParams{"ldap", "tcp", "example.com"},
		},
		{ // with "srv+", "ldap" service and specify port
			"tcp",
			"srv+ldap+example.com:443",
			dialContextResults{&testConn{"tcp", "ldap.example.com:1234"}, nil},
			&dialContextParams{"tcp", "ldap.example.com:1234"},
			&lookupSRVParams{"ldap", "tcp", "example.com"},
		},
	} {
		driver := &testDriver{}
		d := Dialer{drv: driver}
		gotConn, gotErr := d.Dial(c.network, c.address)
		if gotErr != nil {
			if c.want.err == nil {
				t.Errorf("unexpected error: %v", gotErr)
				continue
			}
			if want, got := c.want.err.Error(), gotErr.Error(); got != want {
				t.Errorf("unmatch error:\nwant=%v\ngot=%v", want, got)
			}
			continue
		}
		if want, got := c.want.conn, gotConn; !reflect.DeepEqual(want, got) {
			t.Errorf("unmatch conn:\nwant=%+v\ngot=%+v", want, got)
		}
		if want, got := c.wantDialParams, driver.dialContextParams; !reflect.DeepEqual(want, got) {
			t.Errorf("unmatch dial params:\nwant=%+v\ngot=%+v", want, got)
		}
		if want, got := c.wantLookupParams, driver.lookupSRVParams; !reflect.DeepEqual(want, got) {
			t.Errorf("unmatch lookup params:\nwant=%+v\ngot=%+v", want, got)
		}
	}
}

func TestNew(t *testing.T) {
	d := New(nil)
	if want, got := (&net.Dialer{}), d.drv.(*netDialerDriver).Dialer; !reflect.DeepEqual(want, got) {
		t.Errorf("unexpected underlying dialer:\nwant=%+v\ngot=%+v", want, got)
	}
}
