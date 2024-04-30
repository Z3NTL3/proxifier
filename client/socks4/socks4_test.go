package socks4_test

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"testing"
	"time"

	"github.com/z3ntl3/socks/client"
	"github.com/z3ntl3/socks/client/socks4"
)

func TestSOCKS4TLSClient(t *testing.T) {
	target := client.Context{
		Resolver: net.ParseIP("34.196.110.25"),
		Port:     443,
	}

	proxy := client.Context{
		Resolver: net.ParseIP("72.206.181.97"),
		Port:     64943,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	client, err := socks4.New(target, proxy)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Connect(socks4.USER_NULL, ctx); err != nil {
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
SOCKS4 (TLS):
❯ go test -timeout 30s -run ^TestSOCKS4Client$ github.com/z3ntl3/socks/client/socks4 -v
=== RUN   TestSOCKS4Client
    socks4_test.go:53: HTTP/1.1 200 OK
        Date: Tue, 30 Apr 2024 18:03:27 GMT
        Content-Type: application/json
        Content-Length: 32
        Connection: close
        Server: gunicorn/19.9.0
        Access-Control-Allow-Origin: *
        Access-Control-Allow-Credentials: true

        {
          "origin": "72.206.181.97"
        }

--- PASS: TestSOCKS4Client (10.49s)
*/
