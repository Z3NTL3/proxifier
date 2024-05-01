package client

import (
	"net"
)

type (
	Context struct {
		Resolver interface{}
		Port     int
	}

	Client struct {
		*net.TCPConn
		Target Context
		Proxy  Context
		Worker chan error
	}

	Socks4Client struct{}
	Socks5Client struct {
		Auth
	}

	Auth struct {
		Username string
		Password string
	}

	SocksClient interface {
		*Socks4Client | *Socks5Client
	}
)

/*
SOCKS version 4/4a/5 client
*/
func New[T SocksClient](target Context, proxy Context) T {
	var a interface{}
	b, _ := a.(T)

	return b
}
