# koron-go/dialsrv

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron-go/dialsrv)](https://pkg.go.dev/github.com/koron-go/dialsrv)
[![Actions/Go](https://github.com/koron-go/dialsrv/workflows/Go/badge.svg)](https://github.com/koron-go/dialsrv/actions?query=workflow%3AGo)

Dialer with SRV lookup.

## Sample codes

### HTTP

```go
import "github.com/koron-go/dialsrv"

// query SRV records for "_myservice._tcp.example.com" then make HTTP GET.
r, err := dialsrv.HTTPClient.Get("http://srv+myservice+example.com/")
// TODO: work with r and err as usual.
```

```go
// GET from example.com without SRV.
r, err := dialsrv.HTTPClient.Get("http://example.com/")
// TODO: work with r and err as usual.
```

```go
import "net/http"

// replace http.DefaultClient
http.DefaultClient = dialsrv.HTTPClient
r, err := http.Get("http://srv+myservice+example.com/")
// TODO: work with r and err as usual.
```

### Dialer

```go
import "net"
import "github.com/koron-go/dialsrv"

// wrap *net.Dialer with SRV record querying
d := dialsrv.New(&net.Dialer{/* net.Dialer with your configurations */})
```

## Acceptable formats

*   `srv+{service}+{hostname}` - SRV for `_{service}._{network}.{hostname}`
*   `srv+{hostname}` - SRV for `{hostname}`
