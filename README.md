![1714045920235](image/README/1714045920235.png)

Reliable, multi-tasked and swift SOCKS **connect client**. Implements version ``4/4a/5.``


#### TODO

* [x] TLS support - 
* [x] Version ``4`` 
* [ ] Version ``4/a`` support
* [ ] Version ``5`` 

#### Example SOCKS4 TLS connection stream 
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
		Resolver: net.ParseIP("34.196.110.25"),
		Port:     443,
	}

	proxy := socks.Context{
		Resolver: net.ParseIP("72.206.181.97"),
		Port:     64943,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	client, err := socks.New[*socks.Socks4Client](target, proxy)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Connect(socks.UID_NULL, ctx); err != nil {
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
socks on  main [!?] via 🐹 v1.22.2 took 2s
❯ go run .
HTTP/1.1 200 OK
Date: Wed, 01 May 2024 11:24:44 GMT
Content-Type: application/json
Content-Length: 32
Connection: close
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "origin": "72.206.181.97"
}
*/
```