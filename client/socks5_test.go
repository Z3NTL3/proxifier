package client_test

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	socks "github.com/z3ntl3/socks/client"
)

// go test -timeout 30s -run ^TestSOCKS5Client_NoAuth$ github.com/z3ntl3/socks/client -v
func TestSOCKS5Client_NoAuth(t *testing.T){
	target := socks.Context{
		Resolver: net.ParseIP("149.202.52.226"),
		Port:     80,
	}

	proxy := socks.Context{
		Resolver: net.ParseIP("38.154.227.167"),
		Port:     5868,
	}

	client, err := socks.New(&socks.Socks5Client{},target, proxy)
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


	if _, err := client.Write([]byte("GET / HTTP/1.1\r\nHost: pool.proxyspace.pro\r\nConnection: close\r\n\r\n")); err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(client)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

/*
socks on  main [!] via 🐹 v1.22.2 took 2s
❯ go test -timeout 30s -run ^TestSOCKS5Client_NoAuth$ github.com/z3ntl3/socks/client -v
=== RUN   TestSOCKS5Client_NoAuth
☺&����↨HTTP/1.1 200 OK
Server: nginx/1.18.0 (Ubuntu)
Date: Wed, 01 May 2024 21:37:10 GMT
Content-Type: text/plain
Content-Length: 15
Connection: close

38.154.227.167

--- PASS: TestSOCKS5Client_NoAuth (1.07s)
PASS
ok      github.com/z3ntl3/socks/client  1.249s
*/


// go test -timeout 30s -run ^TestSOCKS5Client_Auth$ github.com/z3ntl3/socks/client -v
func TestSOCKS5Client_Auth(t *testing.T){
	target := socks.Context{
		Resolver: net.ParseIP("149.202.52.226"),
		Port:     80,
	}

	proxy := socks.Context{
		Resolver: net.ParseIP("38.154.227.167"),
		Port:     5868,
	}

	client, err := socks.New(&socks.Socks5Client{},target, proxy)
	if err != nil {
		t.Fatal(err)
	}
	
	{
		client.Auth.Username = "lqafmzlx"
		client.Auth.Password = "i9mzzjv4qdz2"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	
	if err := socks.Connect(client, ctx); err != nil {
		t.Fatal(err)
	}

	defer client.Close()
	client.SetLinger(0)


	if _, err := client.Write([]byte("GET / HTTP/1.1\r\nHost: pool.proxyspace.pro\r\nConnection: close\r\n\r\n")); err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(client)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

/*
socks on  main [!] via 🐹 v1.22.2 took 2s
❯ go test -timeout 30s -run ^TestSOCKS5Client_Auth$ github.com/z3ntl3/socks/client -v
=== RUN   TestSOCKS5Client_Auth
☺&�㧧UHTTP/1.1 200 OK
Server: nginx/1.18.0 (Ubuntu)
Date: Wed, 01 May 2024 21:38:28 GMT
Content-Type: text/plain
Content-Length: 15
Connection: close

38.154.227.167

--- PASS: TestSOCKS5Client_Auth (1.35s)
PASS
ok      github.com/z3ntl3/socks/client  1.537s
*/