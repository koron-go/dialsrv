package dialsrv

import (
	"net"
	"net/http"
	"time"
)

// HTTPTransport is replacement for http.DefaultTransport
var HTTPTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: New(&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// HTTPClient is replacement for http.DefaultClient
var HTTPClient = &http.Client{
	Transport: HTTPTransport,
}
