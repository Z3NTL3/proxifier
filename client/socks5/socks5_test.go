package socks5_test

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/z3ntl3/socks/client"
	"github.com/z3ntl3/socks/client/socks5"
)

func TestSOCKS5Client(t *testing.T) {
	target := client.Context{
		Resolver: net.ParseIP("34.196.110.25"),
		Port:     80,
	}

	proxy := client.Context{
		Resolver: net.ParseIP("38.154.227.167"),
		Port:     5868,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	client, err := socks5.New(target, proxy)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	defer client.Close()
	client.SetLinger(0)

	if _, err := client.Write([]byte("GET /ip HTTP/1.1\r\nHost: httpbin.org\r\nConnection: close\r\n\r\n")); err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(client)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

/*
socks on  main [!?] via 🐹 v1.22.2 took 11s
❯ go test -timeout 30s -run ^TestSOCKS5Client$ github.com/z3ntl3/socks/client/socks5 -v
=== RUN   TestSOCKS5Client
    socks5_test.go:48: ☺&�㧚+HTTP/1.1 200 OK
        Date: Tue, 30 Apr 2024 18:04:55 GMT
        Content-Type: application/json
        Content-Length: 33
        Connection: close
        Server: gunicorn/19.9.0
        Access-Control-Allow-Origin: *
        Access-Control-Allow-Credentials: true

        {
          "origin": "38.154.227.167"
        }

--- PASS: TestSOCKS5Client (2.95s)
PASS
ok      github.com/z3ntl3/socks/client/socks5   3.092s

*/
