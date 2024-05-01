package client_test

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	socks "github.com/z3ntl3/socks/client"
)

// go test -timeout 30s -run ^TestSOCKS5Client_NoAUTH$ github.com/z3ntl3/socks/client -v
func TestSOCKS5Client_NoAUTH(t *testing.T) {
	target := socks.Context{
		Resolver: net.ParseIP("34.196.110.25"),
		Port:     80,
	}

	proxy := socks.Context{
		Resolver: net.ParseIP("38.154.227.167"),
		Port:     5868,
	}

	client, err := socks.New[*socks.Socks5Client](target, proxy)
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
socks on  main [!?] via 🐹 v1.22.2 took 3s
❯ go test -timeout 30s -run ^TestSOCKS5Client_NoAUTH$ github.com/z3ntl3/socks/client -v
=== RUN   TestSOCKS5Client_NoAUTH
    socks5_test.go:52: ☺&����YHTTP/1.1 200 OK
        Date: Wed, 01 May 2024 14:04:14 GMT
        Content-Type: application/json
        Content-Length: 33
        Connection: close
        Server: gunicorn/19.9.0
        Access-Control-Allow-Origin: *
        Access-Control-Allow-Credentials: true

        {
          "origin": "38.154.227.167"
        }

--- PASS: TestSOCKS5Client_NoAUTH (2.28s)
PASS
ok      github.com/z3ntl3/socks/client  2.468s
*/

// go  test -timeout 30s -run ^TestSOCKS5Client_Auth$ github.com/z3ntl3/socks/client -v
func TestSOCKS5Client_Auth(t *testing.T) {
	target := socks.Context{
		Resolver: net.ParseIP("34.196.110.25"),
		Port:     80,
	}

	proxy := socks.Context{
		Resolver: net.ParseIP("38.154.227.167"),
		Port:     5868,
	}

	client, err := socks.New[*socks.Socks5Client](target, proxy)
	if err != nil {
		t.Fatal(err)
	}
	client.Username = "lqafmzlx"
	client.Password = "i9mzzjv4qdz2"


	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := socks.Connect(client, ctx); err != nil {
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
socks on  main [!?] via 🐹 v1.22.2 
❯ go  test -timeout 30s -run ^TestSOCKS5Client_Auth$ github.com/z3ntl3/socks/client -v
=== RUN   TestSOCKS5Client_Auth
    socks5_test.go:113: ☺&���ΏHTTP/1.1 200 OK
        Date: Wed, 01 May 2024 14:26:52 GMT
        Content-Type: application/json
        Content-Length: 33
        Connection: close
        Server: gunicorn/19.9.0
        Access-Control-Allow-Origin: *
        Access-Control-Allow-Credentials: true

        {
          "origin": "38.154.227.167"
        }

--- PASS: TestSOCKS5Client_Auth (1.65s)
PASS
ok      github.com/z3ntl3/socks/client  1.973s

*/