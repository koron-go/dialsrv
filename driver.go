package dialsrv

import (
	"context"
	"net"
)

type driver interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
	LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error)
}

type netDialerDriver struct {
	*net.Dialer
}

var _ driver = (*netDialerDriver)(nil)

func (ndd *netDialerDriver) LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
	r := ndd.Resolver
	if r == nil {
		r = net.DefaultResolver
	}
	return r.LookupSRV(ctx, service, proto, name)
}
