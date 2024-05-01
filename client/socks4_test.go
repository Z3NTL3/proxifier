package client_test

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"testing"
	"time"

	socks "github.com/z3ntl3/socks/client"
)

// go test -timeout 30s -run ^TestSOCKS4Client$ github.com/z3ntl3/socks/client -v
func TestSOCKS4Client(t *testing.T) {
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
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	
	if err := socks.Connect(client, ctx); err != nil {
		t.Fatal(err)
	}

	defer client.Close()
	client.SetLinger(0)

	tlsConn := tls.Client(client, &tls.Config{
		InsecureSkipVerify: true,
	})

	if _, err := tlsConn.Write([]byte("GET /ip HTTP/1.1\r\nHost: httpbin.org\r\nConnection: close\r\n\r\n")); err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(tlsConn)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

/*
socks on  main [!] via 🐹 v1.22.2 
❯ go test -timeout 30s -run ^TestSOCKS4Client$ github.com/z3ntl3/socks/client -v
=== RUN   TestSOCKS4Client
    socks4_test.go:52: HTTP/1.1 200 OK
        Date: Wed, 01 May 2024 11:23:02 GMT
        Content-Type: application/json
        Content-Length: 32
        Connection: close
        Server: gunicorn/19.9.0
        Access-Control-Allow-Origin: *
        Access-Control-Allow-Credentials: true

        {
          "origin": "72.206.181.97"
        }

--- PASS: TestSOCKS4Client (1.18s)
PASS
ok      github.com/z3ntl3/socks/client  1.362s
*/