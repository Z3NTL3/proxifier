package proxifier_test

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/Z3NTL3/proxifier"
)

// go test -timeout 30s -run ^TestHTTP$ github.com/Z3NTL3/proxifier/tests -v
func TestHTTP(t *testing.T) {
	httpClient := proxifier.HTTPClient{}

	conn, err := httpClient.PROXY("https://httpbin.org/ip", proxifier.Context{
		Resolver: net.ParseIP("85.209.2.126"),
		Port: 4444,
	}, time.Second * 10); if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	resp, err := io.ReadAll(conn)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(resp))
}

/*
socks on  main via 🐹 v1.22.2 
❯ go test -timeout 30s -run ^TestHTTP$ github.com/Z3NTL3/proxifier/tests -v
=== RUN   TestHTTP
    http_test.go:30: HTTP/1.1 200 OK
        Server: nginx/1.14.0 (Ubuntu)
        Date: Fri, 17 May 2024 15:29:44 GMT
        Content-Type: application/json
        Content-Length: 31
        Connection: close
        Access-Control-Allow-Origin: *
        Access-Control-Allow-Credentials: true

        {
          "origin": "85.209.2.126"
        }

--- PASS: TestHTTP (1.92s)
PASS
ok      github.com/Z3NTL3/proxifier/tests       2.095s
*/

// go test -timeout 30s -run ^TestHTTPS$ github.com/Z3NTL3/proxifier/tests -v
func TestHTTPS(t *testing.T) {
	httpClient := proxifier.HTTPClient{
		TLS: true,
	}

	conn, err := httpClient.PROXY("https://httpbin.org/ip", proxifier.Context{
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
socks on  main via 🐹 v1.22.2 took 7s
❯ go test -timeout 30s -run ^TestHTTPS$ github.com/Z3NTL3/proxifier/tests -v
=== RUN   TestHTTPS
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Server: Payara Server  4.1.1.164 #badassfish
Date: Fri, 17 May 2024 15:30:27 GMT
Content-Length: 32
Connection: close

{
  "origin": "117.74.65.207"
}

--- PASS: TestHTTPS (5.98s)
PASS
ok      github.com/Z3NTL3/proxifier/tests       6.163s
*/