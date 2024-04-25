![1714045920235](image/README/1714045920235.png)

Reliable, multi-tasked and swift SOCKS connect client. Implements version ``4/4a/5.``

Used as production at <a href="https://pro.simpaix.net">**Lightning CLI </a>**

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
	"log"
	"time"

	"github.com/z3ntl3/socks/client"
	"github.com/z3ntl3/socks/client/socks4"
)

func main() {
	target := client.Context{
		IP:   "172.217.16.206",
		Port: 443,
	}

	proxy := client.Context{
		IP:   "45.81.232.17",
		Port: 12873,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	client := socks4.New(target, proxy)
	if err := client.Connect([]byte{0x00}, ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Close()
    client.SetLinger(0)

	tlsConn := tls.Client(client, &tls.Config{
		InsecureSkipVerify: true,
	})

	buff := make([]byte, 1042) // read first few headers only

	if _, err := tlsConn.Write(
		[]byte("GET / HTTP/1.1\r\nHost: google.com\r\nConnection: close\r\n\r\n"),
	); err != nil {
		log.Fatal(err)
	}
	if _, err := tlsConn.Read(buff); err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(buff))
}

/*
❯ go run .
HTTP/1.1 302 Found
Location: https://www.google.com/sorry/index?continue=https://google.com/&q=EgQtUegRGMeFqbEGIjBU_ngbO4NphFVjEZVaiHD0Vx87jAA_4vXkOHCQ6Rj8HwcVz3qI0sks_-NufuJeV5EyAXJaAUM
x-hallmonitor-challenge: CgsIyIWpsQYQo9yUEBIELVHoEQ
Content-Type: text/html; charset=UTF-8
Content-Security-Policy-Report-Only: object-src 'none';base-uri 'self';script-src 'nonce-K5LptrfOVe6LDj7xGGqOoA' 'strict-dynamic' 'report-sample' 'unsafe-eval' 'unsafe-inline' https: http:;report-uri https://csp.withgoogle.com/csp/gws/other-hp
P3P: CP="This is not a P3P policy! See g.co/p3phelp for more info."
Date: Thu, 25 Apr 2024 11:47:20 GMT
Server: gws
Content-Length: 358
X-XSS-Protection: 0
X-Frame-Options: SAMEORIGIN
Set-Cookie: AEC=AQTF6Hzw6A6ZGFOvhD3drLEnFNyCI1dNJP9JHdaJDYhhRxi9TQzpSdZFKtc; expires=Tue, 22-Oct-2024 11:47:20 GMT; path=/; domain=.google.com; Secure; HttpOnly; SameSite=lax
Set-Cookie: __Secure-ENID=19.SE=U-7cnG-S9fxFKtZu3R4-0LGlPETmvWv_6bWCLbNy0_veLAXQaSb_HSzDzHyB1kZLs2fO1SfYJFzppeVWaeoghIjZHm_FdZNJ3o3IZU-0tP7s-MIypoHzA
*/

```
