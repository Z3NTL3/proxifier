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

func TestSocks4Client(t *testing.T) {
	target := client.Context{
		Resolver: net.ParseIP("149.202.52.226"),
		Port:     443,
	}

	proxy := client.Context{
		Resolver: net.ParseIP("192.162.232.15"),
		Port:     1080,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	client, err := socks4.New(target, proxy)
	if err != nil {
		t.Fatal(err)
	}

	// socks4.NULL means no user-id.
	if err := client.Connect([]byte{socks4.NULL}, ctx); err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	client.SetLinger(0)

	tlsConn := tls.Client(client, &tls.Config{
		InsecureSkipVerify: true,
	})

	if _, err := tlsConn.Write([]byte("GET /ip HTTP/1.1\r\nHost: pool.proxyspace.pro\r\nConnection: close\r\n\r\n")); err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(tlsConn)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

/*
❯ go test -timeout 30s -run ^TestSocks4Client$ github.com/z3ntl3/socks/client/socks4 -v
=== RUN   TestSocks4Client
    socks4_test.go:54: HTTP/1.1 200 OK
        Server: nginx/1.18.0 (Ubuntu)
        Date: Mon, 29 Apr 2024 18:41:48 GMT
        Content-Type: text/plain
        Content-Length: 15
        Connection: close

        192.162.232.15

--- PASS: TestSocks4Client (0.61s)
PASS
ok      github.com/z3ntl3/socks/client/socks4   0.778s
*/
