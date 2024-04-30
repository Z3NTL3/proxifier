/*
Socks4 connect implementation
*/
package socks4

import (
	"context"
	"encoding/binary"
	"errors"
	"net"

	"github.com/z3ntl3/socks/client"
)

type (
	reply = byte // convenience
)

// no export
// fix: unnecessary types
const (
	version byte = 0x04
	cmd     byte = 0x01
	null    byte = 0x00

	granted reply = 0x5A // Request granted
)

var (
	USER_NULL  []byte = []byte{0x00} // for convenience
	reply_enum        = map[reply]string{
		0x5A: "Request granted",
		0x5B: "Request rejected or failed",
		0x5C: "Request failed because client is not running identd (or not reachable from server)",
		0x5D: "Request failed because client's identd could not confirm the user ID in the request",
	}
)

type Socks4Client struct {
	client.Client
}

/*
Creates a new SOCKS4 Connect client
*/
func New(target client.Context, proxy client.Context) (client_ *Socks4Client, err error) {
	// godoc states:
	// Deferred functions may read and assign to the returning function’s named return values.
	defer func() {
		panicErr := recover()
		if panicErr != nil {
			err = panicErr.(error)
		}
	}()

	// socks4 client requires net.ip
	resolvers := []net.IP{target.Resolver.(net.IP), proxy.Resolver.(net.IP)}

	if valid := client.IsIPV4(resolvers[0], resolvers[1]); !valid {
		return nil, err
	}

	client_ = new(Socks4Client)
	{
		client_.Client = client.Client{
			Target: target,
			Proxy:  proxy,
			Worker: make(chan error),
		}
	}

	return
}

/*
Connects to the target through proxy and returns proxy tunnel
*/
func (c *Socks4Client) Connect(uid []byte, ctx context.Context) error {
	has_null := false // has termination byte ? (required)
	for _, b := range uid {
		if b == null {
			has_null = true
		}
	}

	if !has_null {
		uid = append(uid, null)
	}

	go func() {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   c.Proxy.Resolver.(net.IP),
			Port: c.Proxy.Port,
		})
		if err != nil {
			c.Worker <- err
			return
		}

		c.TCPConn = conn
		go c.tunnel(uid)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c.Worker:
		close(c.Worker)
		return err
	}
}

func (c *Socks4Client) tunnel(uid []byte) {
	var err error
	// shallow copy
	defer func(sh_clone *Socks4Client, err_ *error) {
		if *err_ != nil {
			sh_clone.Close()
		}

		panicErr := recover()
		if panicErr != nil {
			*err_ = errors.New(panicErr.(string))
		}

		sh_clone.Worker <- *err_
	}(c, &err)

	var HEADER []byte
	{
		PORT := make([]byte, 2)
		binary.BigEndian.PutUint16(PORT, uint16(c.Target.Port))

		HEADER = append(HEADER, version)
		HEADER = append(HEADER, cmd)
		HEADER = append(HEADER, PORT...)
		HEADER = append(HEADER, c.Target.Resolver.(net.IP).To4()...)
		HEADER = append(HEADER, uid...)

	}

	n, err := c.Write(HEADER)
	if err != nil || !(n > 0) {
		if !(n > 0) {
			err = client.ErrFailedHeaderPacket
		}
		return
	}

	HEADER = make([]byte, 8)
	if _, err = c.Read(HEADER); err != nil {
		return
	}

	if HEADER[1] != granted {
		err = errors.New(reply_enum[HEADER[1]])
	}
}
