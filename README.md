![1714045920235](image/README/1714045920235.png)

Reliable, multi-tasked and swift SOCKS **connect client**. Implements version ``4/4a/5.``


#### TODO

* [x] TLS support - 
* [x] Version ``4`` 
* [ ] Version ``4/a`` support
* [x] Version ``5`` 

#### Example SOCKS4 TLS connection stream 
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