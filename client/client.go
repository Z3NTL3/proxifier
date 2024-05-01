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
		target Context
		proxy  Context
		worker chan error
	}

	Socks4Client struct {
		UID []byte
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
		setup()
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

	client = *new(T)

	props := Client{
		target: target,
		proxy:  proxy,
		worker: make(chan error),
	}

	switch any(client).(type) {

		case *Socks4Client:
			// not valid ipv 4
			if !IsIPV4(target.Resolver.(net.IP), proxy.Resolver.(net.IP)) { // may panic
				err = ErrNotValidIP
			}
			client = any(&Socks4Client{
				Client: props,
			}).(T)
			
		case *Socks5Client:
			// not ipv 4,6 or is not a host
			if !IsAccepted(target.Resolver, proxy.Resolver) {
				err = ErrUnsupported
			}

			domain, ok := props.target.Resolver.(string)
			if ok {
				if !IsDomain(domain) {
					err = ErrNotValidDomain
				}
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

/*
	Tunnels through proxy to target. On failure returns error
*/
func Connect[T SocksClient](client T,  ctx context.Context) error {
	var worker chan error

	client.setup()
	
	switch c := any(client).(type) {
		case *Socks4Client:
			worker = c.worker
		case *Socks5Client:
			if len(c.Username) > 0 && len(c.Password) > 0 {
				if !MinChar(c.Username, c.Password) {
					return ErrMax255Char
				}
			}
			worker = c.worker
		default:
			return ErrNotValidClient
	}

	select {
		case <-ctx.Done():
			return ctx.Err()

		case err := <-worker:
			close(worker)
			return err
	}
}