package client

import (
	"context"
	"net"
)

type Context struct {
	Resolver interface{}
	Port     int
}

type Client struct {
	*net.TCPConn
	target *Context
	proxy  *Context
	worker chan error
}

type Socks4Client struct {
	Client
	UID []byte
}

type Socks5Client struct {
	Client
	Auth
}

type Auth struct {
	Username string
	Password string
}

type SocksClient interface {
	*Socks4Client | *Socks5Client
	setup() chan error
	init(target,proxy *Context) error
}

type Command = byte

const (
	CMD Command = 0x01 // stream connection
	NULL byte = 0x00
)

/*
SOCKS version 4/4a/5 client

currently only v4/5
*/
func New[T SocksClient](client T, target, proxy Context)(T, error) {
	return client, client.init(&target, &proxy)
}

/*
	Tunnels through proxy to target. On failure returns error
*/
func Connect[T SocksClient](client T,  ctx context.Context) error {
	worker := client.setup()

	select {
		case <-ctx.Done():
			return ctx.Err()

		case err := <-worker:
			close(worker)
			return err
	}
}