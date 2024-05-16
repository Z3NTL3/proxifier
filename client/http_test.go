package client_test

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	socks "github.com/z3ntl3/socks/client"
)

// go test -timeout 30s -run ^TestHTTP$ github.com/z3ntl3/socks/client -v
func TestHTTP(t *testing.T) {
	httpClient := socks.HTTPClient{}

	conn, err := httpClient.PROXY("https://httpbin.org/ip", socks.Context{
		Resolver: net.ParseIP("103.153.142.18"),
		Port: 8080,
	}, time.Second * 5); if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	resp, err := io.ReadAll(conn)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(resp))
}

/*
socks on  main [!?] via 🐹 v1.22.2 
❯ go test -timeout 30s -run ^TestHTTP$ github.com/z3ntl3/socks/client -v
=== RUN   TestHTTP
HTTP/1.1 200 OK
Server: nginx/1.14.0 (Ubuntu)
Date: Thu, 16 May 2024 17:23:02 GMT
Content-Type: application/json
Content-Length: 34
Connection: close
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "origin": "193.162.143.137"
}

--- PASS: TestHTTP (0.70s)
PASS
ok      github.com/z3ntl3/socks/client  0.912s
*/

// go test -timeout 30s -run ^TestHTTPS$ github.com/z3ntl3/socks/client -v
func TestHTTPS(t *testing.T) {
	httpClient := socks.HTTPClient{
		TLS: true,
	}

	conn, err := httpClient.PROXY("https://httpbin.org/ip", socks.Context{
		Resolver: net.ParseIP("117.74.65.207"),
		Port: 54417,
	}, time.Second * 10); if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
	})

	if _, err = tlsConn.Write([]byte("GET /ip HTTP/1.1\r\nHost: httpbin.org\r\nConnection: close\r\n\r\n")); err != nil {
		t.Fatal(err)
	}
	
	resp, err := io.ReadAll(tlsConn)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(resp))
}

/*
socks on  main [!?] via 🐹 v1.22.2 took 5s
❯ go test -timeout 30s -run ^TestHTTPS$ github.com/z3ntl3/socks/client -v
=== RUN   TestHTTPS

HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Server: gSOAP/2.8
Date: Thu, 16 May 2024 17:35:55 GMT
Content-Length: 32
Connection: close

{
  "origin": "117.74.65.207"
}

--- PASS: TestHTTPS (6.05s)
PASS
ok      github.com/z3ntl3/socks/client  6.235s
*/