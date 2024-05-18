package proxifier_test

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"testing"
	"time"

	"github.com/Z3NTL3/proxifier"
)

// go test -timeout 30s -run ^TestSOCKS5Client_NoAuth$ github.com/Z3NTL3/proxifier/tests -v
func TestSOCKS5Client_NoAuth(t *testing.T){
	addr, err := proxifier.LookupHost("httpbin.org")
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	
	if err := proxifier.Connect(client, ctx); err != nil {
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
❯ go test -timeout 30s -run ^TestSOCKS5Client_NoAuth$ github.com/Z3NTL3/proxifier/tests -v
=== RUN   TestSOCKS5Client_NoAuth
    socks5_test.go:60: HTTP/1.1 200 OK
        Date: Fri, 17 May 2024 12:22:31 GMT
        Content-Type: application/json
        Content-Length: 33
        Connection: close
        Server: gunicorn/19.9.0
        Access-Control-Allow-Origin: *
        Access-Control-Allow-Credentials: true

        {
          "origin": "38.154.227.167"
        }

--- PASS: TestSOCKS5Client_NoAuth (1.32s)
PASS
ok      github.com/Z3NTL3/proxifier/tests      1.644s
*/


// go test -timeout 30s -run ^TestSOCKS5Client_Auth$ github.com/Z3NTL3/proxifier/tests -v
func TestSOCKS5Client_Auth(t *testing.T){
	addr, err := proxifier.LookupHost("httpbin.org")
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}
	
	{
		client.Auth.Username = "lqafmzlx"
		client.Auth.Password = "i9mzzjv4qdz2"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	
	if err := proxifier.Connect(client, ctx); err != nil {
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
socks on  main [!] via 🐹 v1.22.2 took 2s
❯ go test -timeout 30s -run ^TestSOCKS5Client_Auth$ github.com/Z3NTL3/proxifier/tests -v
=== RUN   TestSOCKS5Client_Auth
    socks5_test.go:115: HTTP/1.1 200 OK
        Date: Fri, 17 May 2024 12:21:11 GMT
        Content-Type: application/json
        Content-Length: 33
        Connection: close
        Server: gunicorn/19.9.0
        Access-Control-Allow-Origin: *
        Access-Control-Allow-Credentials: true

        {
          "origin": "38.154.227.167"
        }

--- PASS: TestSOCKS5Client_Auth (1.27s)
PASS
ok      github.com/Z3NTL3/proxifier/tests      1.584s
*/