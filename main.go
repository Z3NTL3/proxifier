package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/z3ntl3/socks/client"
	"github.com/z3ntl3/socks/client/socks4"
	"github.com/z3ntl3/socks/client/socks5"
)

func main() {
	// test()
	// time.Sleep(time.Second * 1)

	target := client.Context{
		Resolver: net.ParseIP("149.202.52.226"),
		Port:     443,
	}

	proxy := client.Context{
		Resolver: net.ParseIP("192.252.208.67"),
		Port:     14287,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	client, err := socks5.New(target, proxy)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	client.SetLinger(0)

	tlsConn := tls.Client(client, &tls.Config{
		InsecureSkipVerify: true,
	})

	if _, err := tlsConn.Write([]byte("GET / HTTP/1.1\r\nHost: pool.proxyspace.pro\r\nConnection: close\r\n\r\n")); err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(tlsConn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}

func test() {
	target := client.Context{
		Resolver: net.ParseIP("149.202.52.226"),
		Port:     443,
	}

	proxy := client.Context{
		Resolver: net.ParseIP("199.229.254.129"),
		Port:     4145,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	client, err := socks4.New(target, proxy)
	if err != nil {
		log.Fatal(err)
	}

	// socks4.NULL means no user-id.
	if err := client.Connect(socks4.USER_NULL, ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	client.SetLinger(0)

	tlsConn := tls.Client(client, &tls.Config{
		InsecureSkipVerify: true,
	})

	if _, err := tlsConn.Write([]byte("GET /ip HTTP/1.1\r\nHost: pool.proxyspace.pro\r\nConnection: close\r\n\r\n")); err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(tlsConn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}
