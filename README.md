![1714045920235](image/README/1714045920235.png)

Reliable, multi-tasked and swift SOCKS **connect client**. Implements version ``4/4a/5.``


#### TODO

* [x] TLS support 
* [x] Version ``4`` 
* [ ] Version ``4/a`` support
* [x] Version ``5`` 

## Examples
You can find more in files that include ``_test``

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

	client, err := socks.New[*socks.Socks4Client](target, proxy)
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
socks on  main via 🐹 v1.22.2 took 2s
❯ go run main.go
HTTP/1.1 200 OK
Date: Wed, 01 May 2024 14:42:56 GMT
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


#### Example SOCKS5 NO AUTH
```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	socks "github.com/z3ntl3/socks/client"
)

func main() {
	target := socks.Context{
		Resolver: net.ParseIP("3.211.223.136"),
		Port:     80,
	}

	proxy := socks.Context{
		Resolver: net.ParseIP("38.154.227.167"),
		Port:     5868,
	}

	client, err := socks.New[*socks.Socks5Client](target, proxy)
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

	if _, err := client.Write([]byte("GET /ip HTTP/1.1\r\nHost: httpbin.org\r\nConnection: close\r\n\r\n")); err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}

/*

{
  "origin": "38.154.227.167"
}
*/

```