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

	Socks4Client struct {
		Client
	}
	Socks5Client struct {
		Auth
		Client
	}

	Auth struct {
		Username string
		Password string
	}

	SocksClient interface {
		*Socks4Client | *Socks5Client
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
func New[T SocksClient](target Context, proxy Context) (client T, err error) {
	// according to go doc
	// defer func may assign to named returns
	defer func() {
		panicErr := recover()
		if panicErr != nil {
			err = panicErr.(error)
			// may also panic but i know its always an error type so
		}
	}()

	props := Client{
		Target: target,
		Proxy:  proxy,
		Worker: make(chan error),
	}

	switch any(client).(type) {

		case *Socks4Client:
			// not valid ipv 4
			if !IsIPV4(target.Resolver.(net.IP), proxy.Resolver.(net.IP)) { // may panic
				err = ErrNotValidIP
			}
			client = any(&Socks4Client{
				props,
			}).(T)
			
		case *Socks5Client:
			// not ipv 4,6 or is not a host
			if !IsAccepted(target.Resolver, proxy.Resolver) {
				err = ErrUnsupported
			}
			client = any(&Socks5Client{
				Client: props,
				Auth:   Auth{}, // public properties
			}).(T)
		default:
			err = ErrNotValidClient
	}

	return
}
