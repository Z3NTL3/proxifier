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

// go test -timeout 30s -run ^TestSOCKS4Client$ github.com/Z3NTL3/proxifier/tests -v
func TestSOCKS4Client(t *testing.T){
	addr, err := proxifier.LookupHost("pool.proxyspace.pro")
	if err != nil {
		t.Fatal(err)
	}

	target := proxifier.Context{
		Resolver: net.ParseIP(addr[0]),
		Port:     443,
	}

	proxy := proxifier.Context{
		Resolver: net.ParseIP("174.64.199.82"),
		Port:     4145,
	}

	client, err := proxifier.New(&proxifier.Socks4Client{},target, proxy)
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

	if _, err := tlsConn.Write(
		[]byte("GET / HTTP/1.1\r\nHost: pool.proxyspace.pro\r\nConnection: close\r\n\r\n"),
	); err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(tlsConn)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}
/*
socks on  main [!] via 🐹 v1.22.2 took 3s
❯ go test -timeout 30s -run ^TestSOCKS4Client$ github.com/Z3NTL3/proxifier/tests -v
=== RUN   TestSOCKS4Client
    socks4_test.go:56: HTTP/1.1 200 OK
        Server: nginx/1.18.0 (Ubuntu)
        Date: Fri, 17 May 2024 12:18:57 GMT
        Content-Type: text/plain
        Content-Length: 14
        Connection: close

        174.64.199.82

--- PASS: TestSOCKS4Client (1.53s)
PASS
ok      github.com/Z3NTL3/proxifier/tests      1.863s
*/
