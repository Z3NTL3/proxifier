package proxifier

import (
	"context"
	"net"
)

// Detail about the network to handshake
type Context struct {
	Resolver interface{}
	Port     int
}

// Core proxy client
type Client struct {
	*net.TCPConn // underyling tcp connection
	target *Context // target context
	proxy  *Context // proxy context
	worker chan error // synchronization primitive
}

// SOCKS5 client
type Socks4Client struct {
	Client // core client
	UID []byte // userid, defaults to null
}

// SOCKS5 client
type Socks5Client struct {
	Client // core client
	Auth // authentication context
}

// Authentication credentials
type Auth struct {
	Username string // Username
	Password string // Password
}

// Proxy client which implements SOCKS4/SOCKS5
type SocksClient interface {
	*Socks4Client | *Socks5Client
	setup() chan error
	init(target, proxy *Context) error
}

// TCP/IP stream
type Command = byte

const (
	CMD  Command = 0x01 // stream connection
	NULL byte    = 0x00 // null byte
)

// Creates new [SocksClient].
//
// ``target`` and ``proxy`` comfort [Context]
//
// On failure, returns error.
func New[T SocksClient](client T, target, proxy Context) (T, error) {
	return client, client.init(&target, &proxy)
}

// Tunnels through proxy to target. On failure returns error.
//
// ``ctx`` is a [context.Context] used for timeout/cancellation signal.
func Connect[T SocksClient](client T, ctx context.Context) error {
	worker := client.setup()

	select {
		case <-ctx.Done():
			return ctx.Err()

		case err := <-worker:
			close(worker)
			return err
	}
}