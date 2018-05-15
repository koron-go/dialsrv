# koron-go/dialsrv

[![GoDoc](https://godoc.org/github.com/koron-go/dialsrv?status.svg)](https://godoc.org/github.com/koron-go/dialsrv)
[![CircleCI](https://img.shields.io/circleci/project/github/koron-go/dialsrv/master.svg)](https://circleci.com/gh/koron-go/dialsrv/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/dialsrv)](https://goreportcard.com/report/github.com/koron-go/dialsrv)

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
