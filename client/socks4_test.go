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
func TestSOCKS4Client(t *testing.T){
	target := socks.Context{
		Resolver: net.ParseIP("149.202.52.226"),
		Port:     443,
	}

	proxy := socks.Context{
		Resolver: net.ParseIP("174.64.199.82"),
		Port:     4145	,
	}

	client, err := socks.New(&socks.Socks4Client{},target, proxy)
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

	if _, err := tlsConn.Write([]byte("GET / HTTP/1.1\r\nHost: pool.proxyspace.pro\r\nConnection: close\r\n\r\n")); err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(tlsConn)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

/*
socks on  main [✘!] via 🐹 v1.22.2 
❯ go test -timeout 30s -run ^TestSOCKS4Client$ github.com/z3ntl3/socks/client -v
=== RUN   TestSOCKS4Client
HTTP/1.1 200 OK
Server: nginx/1.18.0 (Ubuntu)
Date: Wed, 01 May 2024 21:35:09 GMT
Content-Type: text/plain
Content-Length: 14
Connection: close

174.64.199.82

--- PASS: TestSOCKS4Client (2.58s)
PASS
ok      github.com/z3ntl3/socks/client  2.775s
*/