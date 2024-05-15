package client

import (
	"context"
	"net"
)

type (
	Context struct {
		Resolver interface{}
		Port     int
	}

	Client struct {
		*net.TCPConn
		target *Context
		proxy  *Context
		worker chan error
	}

	Socks4Client struct {
		Client
		UID []byte
	}

	Socks5Client struct {
		Client
		Auth
	}

	Auth struct {
		Username string
		Password string
	}

	SocksClient interface {
		*Socks4Client | *Socks5Client
		setup() chan error
		init(target,proxy *Context) error
	}

	Command = byte
)
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