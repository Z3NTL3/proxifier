![1714045920235](image/README/1714045920235.png)

Reliable, multi-tasked and swift SOCKS **connect client**. Implements version ``4/4a/5.``


#### TODO

* [x] TLS support 
* [x] Version ``4`` 
* [ ] Version ``4/a`` support
* [x] Version ``5`` 

## Examples
You may find more examples in ``_test`` files.

#### SOCKS4 TLS
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

	socks "github.com/z3ntl3/socks/client"
)

func main() {
	target := socks.Context{
		Resolver: net.ParseIP("149.202.52.226"),
		Port:     443,
	}

	proxy := socks.Context{
		Resolver: net.ParseIP("174.64.199.82"),
		Port:     4145,
	}

	client, err := socks.New(&socks.Socks4Client{}, target, proxy)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := socks.Connect(client, ctx); err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	client.SetLinger(0)

	tlsConn := tls.Client(client, &tls.Config{
		InsecureSkipVerify: true,
	})

	if _, err := tlsConn.Write([]byte("GET / HTTP/1.1\r\nHost: pool.proxyspace.pro\r\nConnection: close\r\n\r\n")); err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(tlsConn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}

/*
socks on  main [!] via 🐹 v1.22.2 took 2s
❯ go run .
HTTP/1.1 200 OK
Server: nginx/1.18.0 (Ubuntu)
Date: Wed, 01 May 2024 21:39:44 GMT
Content-Type: text/plain
Content-Length: 14
Connection: close

174.64.199.82
*/
```