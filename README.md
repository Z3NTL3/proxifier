![1714045920235](image/README/1714045920235.png)

Reliable, multi-tasked and swift SOCKS **connect client**. Implements version ``4/4a/5.``

#### TODO

* [X] TLS support
* [X] Version ``4`` support
* [ ] Version ``4/a`` support
* [ ] Version ``5`` support

#### Note
- SOCKS4 has been completely implemented for CONNECT
- SOCKS5 implementation currently does only support NO_AUTH
#### Example
##### SOCKS4 (TLS)
```go
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/z3ntl3/socks/client"
	"github.com/z3ntl3/socks/client/socks4"
)

func main() {
	target := client.Context{
		Resolver: net.ParseIP("34.196.110.25"),
		Port:     443,
	}

	proxy := client.Context{
		Resolver: net.ParseIP("190.12.95.170"),
		Port:     37209,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	client, err := socks4.New(target, proxy)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Connect(socks4.USER_NULL, ctx); err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	client.SetLinger(0)

	tlsConn := tls.Client(client, &tls.Config{
		InsecureSkipVerify: true,
	})

	if _, err := tlsConn.Write([]byte("GET /ip HTTP/1.1\r\nHost: httpbin.org\r\nConnection: close\r\n\r\n")); err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(tlsConn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}

/*
SOCKS4 (TLS):
❯ go test -timeout 30s -run ^TestSOCKS4Client$ github.com/z3ntl3/socks/client/socks4 -v
=== RUN   TestSOCKS4Client
    socks4_test.go:53: HTTP/1.1 200 OK
        Date: Tue, 30 Apr 2024 18:03:27 GMT
        Content-Type: application/json
        Content-Length: 32
        Connection: close
        Server: gunicorn/19.9.0
        Access-Control-Allow-Origin: *
        Access-Control-Allow-Credentials: true

        {
          "origin": "72.206.181.97"
        }

--- PASS: TestSOCKS4Client (10.49s)
*/
```