<div align="center">
	<kbd><img src="https://github.com/Z3NTL3/socks/assets/48758770/9c91baa2-03df-480b-8c02-4e082a2f28f5" width="400" border="1px solid red"></kbd>
	<div>
		<img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/z3ntl3/socks">
		<img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/z3ntl3/socks">
		<img alt="Static Badge" src="https://img.shields.io/badge/author-z3ntl3-blue?style=plastic">
	</div>
</div>

# Proxifier

Reliable proxy client library for Go programs.

#### Features

* [X] TLS 
* [X] SOCKS4
* [X] SOCKS5
* [X] HTTP ``(HTTP FORWARD)``
* [X] HTTPS ``(HTTP TUNNEL)``
* [X] Auth

##### Todo
* [ ] SOCKS``4/a`` support


## Examples

#### HTTPS AUTH

HTTPS TUNNEL mechanism with authentication.

```go
package main

import (
	"io"
	"log"
	"net"
	"time"

	"crypto/tls"

	"github.com/Z3NTL3/proxifier"
)

func main() {
	httpClient := proxifier.HTTPClient{
		TLS: true,
		Auth: proxifier.Auth{
			Username: "hello",
			Password: "world",
		},
	}

	conn, err := httpClient.PROXY("https://httpbin.org/ip", proxifier.Context{
		Resolver: net.ParseIP("117.74.65.207"),
		Port: 54417,
	}, time.Second * 10); if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
	})

	if _, err = tlsConn.Write([]byte("GET /ip HTTP/1.1\r\nHost: httpbin.org\r\nConnection: close\r\n\r\n")); err != nil {
		log.Fatal(err)
	}
	
	resp, err := io.ReadAll(tlsConn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(resp))
}
/*
RESULT:

~\Documents\sockstests via 🐹 v1.22.2 took 7s
❯ go run .
2024/05/17 18:49:27 HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Server: CherryPy/3.2.5
Www-Authenticate: Basic realm="Broadband Router"
Date: Fri, 17 May 2024 16:49:28 GMT
Content-Length: 32
Connection: close

{
  "origin": "117.74.65.207"
}
*/
```

#### HTTPS NO AUTH

HTTPS TUNNEL with no authentication

```go
package main

import (
	"io"
	"log"
	"net"
	"time"

	"crypto/tls"

	"github.com/Z3NTL3/proxifier"
)

func main() {
	httpClient := proxifier.HTTPClient{
		TLS: true,
	}

	conn, err := httpClient.PROXY("https://httpbin.org/ip", proxifier.Context{
		Resolver: net.ParseIP("117.74.65.207"),
		Port: 54417,
	}, time.Second * 10); if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
	})

	if _, err = tlsConn.Write([]byte("GET /ip HTTP/1.1\r\nHost: httpbin.org\r\nConnection: close\r\n\r\n")); err != nil {
		log.Fatal(err)
	}
	
	resp, err := io.ReadAll(tlsConn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(resp))
}
/*
RESULT:

~\Documents\sockstests via 🐹 v1.22.2 took 7s
❯ go run .
2024/05/17 18:49:27 HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Server: CherryPy/3.2.5
Www-Authenticate: Basic realm="Broadband Router"
Date: Fri, 17 May 2024 16:49:28 GMT
Content-Length: 32
Connection: close

{
  "origin": "117.74.65.207"
}
*/
```

#### HTTP NO AUTH

HTTP FORWARD with no authentication.

```go
package main

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/Z3NTL3/proxifier"
)

func main() {
	httpClient := proxifier.HTTPClient{}

	conn, err := httpClient.PROXY("https://httpbin.org/ip", proxifier.Context{
		Resolver: net.ParseIP("85.209.2.126"),
		Port:     4444,
	}, time.Second*10)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	resp, err := io.ReadAll(conn)
	if err != nil {
		log.Fatal(err)

	}

	log.Println(string(resp))
}
/*
RESULT:

~\Documents\sockstests via 🐹 v1.22.2 took 3s
❯ go run .
2024/05/17 18:35:41 HTTP/1.1 200 OK
Server: nginx/1.14.0 (Ubuntu)
Date: Fri, 17 May 2024 16:35:41 GMT
Content-Type: application/json
Content-Length: 31
Connection: close
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "origin": "85.209.2.126"
}
*/
```

#### HTTP AUTH
HTTP FORWARD with authentication

```go
package main

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/Z3NTL3/proxifier"
)

func main() {
	httpClient := proxifier.HTTPClient{
		Auth: proxifier.Auth{
			Username: "hello",
			Password: "world",
		},
	}

	conn, err := httpClient.PROXY("https://httpbin.org/ip", proxifier.Context{
		Resolver: net.ParseIP("85.209.2.126"),
		Port:     4444,
	}, time.Second*10)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	resp, err := io.ReadAll(conn)
	if err != nil {
		log.Fatal(err)

	}

	log.Println(string(resp))
}
/*
RESULT:

~\Documents\sockstests via 🐹 v1.22.2 took 3s
❯ go run .
2024/05/17 18:35:41 HTTP/1.1 200 OK
Server: nginx/1.14.0 (Ubuntu)
Date: Fri, 17 May 2024 16:35:41 GMT
Content-Type: application/json
Content-Length: 31
Connection: close
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "origin": "85.209.2.126"
}
*/
```

#### SOCKS5 NO AUTH

With no authentication

```go
package main

import (
	"io"
	"log"
	"net"
	"time"

	"context"
	"crypto/tls"

	"github.com/Z3NTL3/proxifier"
)

func main() {
	addr, err := proxifier.LookupHost("httpbin.org")
	if err != nil {
		log.Fatal(err)
	}

	target := proxifier.Context{
		Resolver: net.ParseIP(addr[0]),
		Port:     443,
	}

	proxy := proxifier.Context{
		Resolver: net.ParseIP("38.154.227.167"),
		Port:     5868,
	}

	client, err := proxifier.New(&proxifier.Socks5Client{},target, proxy)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	
	if err := proxifier.Connect(client, ctx); err != nil {
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

	log.Println(string(data))
}
/*
RESULT:

~\Documents\sockstests via 🐹 v1.22.2 
❯ go run .
2024/05/17 18:44:03 HTTP/1.1 200 OK
Date: Fri, 17 May 2024 16:44:03 GMT
Content-Type: application/json
Content-Length: 33
Connection: close
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "origin": "38.154.227.167"
}
*/
```

#### SOCKS4

```go
package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"time"

	"github.com/Z3NTL3/proxifier"
)

func main() {
	addr, err := proxifier.LookupHost("httpbin.org")
	if err != nil {
		log.Fatal(err)
	}

	target := proxifier.Context{
		Resolver: net.ParseIP(addr[0]),
		Port:     443,
	}

	proxy := proxifier.Context{
		Resolver: net.ParseIP("174.64.199.82"),
		Port:     4145,
	}

	client, err := proxifier.New(&proxifier.Socks4Client{},target, proxy)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	
	if err := proxifier.Connect(client, ctx); err != nil {
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

	log.Println(string(data))
}
/*
RESULT:

~\Documents\sockstests via 🐹 v1.22.2 took 6s
❯ go run .
2024/05/17 18:54:00 HTTP/1.1 200 OK
Date: Fri, 17 May 2024 16:54:00 GMT
Content-Type: application/json
Content-Length: 32
Connection: close
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "origin": "174.64.199.82"
}
*/
```

#### SOCKS5 AUTH

With authentication.

```go
package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"time"

	"github.com/Z3NTL3/proxifier"
)

func main() {
	addr, err := proxifier.LookupHost("httpbin.org")
	if err != nil {
		log.Fatal(err)
	}

	target := proxifier.Context{
		Resolver: net.ParseIP(addr[0]),
		Port:     443,
	}

	proxy := proxifier.Context{
		Resolver: net.ParseIP("38.154.227.167"),
		Port:     5868,
	}

	client, err := proxifier.New(&proxifier.Socks5Client{},target, proxy)
	if err != nil {
		log.Fatal(err)
	}
	
	{
		client.Auth.Username = "lqafmzlx"
		client.Auth.Password = "i9mzzjv4qdz2"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	
	if err := proxifier.Connect(client, ctx); err != nil {
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

	log.Println(string(data))
}
/*
RESULT:

~\Documents\sockstests via 🐹 v1.22.2 
❯ go run .
2024/05/17 18:37:49 HTTP/1.1 200 OK
Date: Fri, 17 May 2024 16:37:49 GMT
Content-Type: application/json
Content-Length: 33
Connection: close
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "origin": "38.154.227.167"
}
*/
```