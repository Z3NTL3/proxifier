![1714045920235](image/README/1714045920235.png)

Reliable, multi-tasked and swift SOCKS connect client. Implements version ``4/4a/5.``

#### TODO

* [X] TLS support
* [X] Version ``4`` support
* [ ] Version ``4/a`` support
* [ ] Version ``5`` support

### Example

##### Socks4 TLS

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
	target := client.TargetCtx{
		IP:   net.ParseIP("149.202.52.226"),
		Port: 443,
	}

	proxy := client.ProxyCtx{
		IP:   net.ParseIP("209.126.7.172"),
		Port: 13637,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	client, err := socks4.New(target, proxy)
	if err != nil {
		log.Fatal(err)
	}

	// socks4.NULL means no user-id.
	if err := client.Connect([]byte{socks4.NULL}, ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	client.SetLinger(0)

	tlsConn := tls.Client(client, &tls.Config{
		InsecureSkipVerify: true,
	})

	if _, err := tlsConn.Write([]byte("GET /ip HTTP/1.1\r\nHost: pool.proxyspace.pro\r\nConnection: close\r\n\r\n")); err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(tlsConn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}

/*
HTTP/1.1 200 OK
Server: nginx/1.18.0 (Ubuntu)
Date: Sat, 27 Apr 2024 18:34:25 GMT
Content-Type: text/plain
Content-Length: 14
Connection: close

209.126.7.172
*/

```
