/*
TODO: 30-04-2024
*/
package socks5

import (
	"context"
	"encoding/binary"
	"errors"
	"net"

	"github.com/z3ntl3/socks/client"
)

type (
	// convenience
	auth      = byte
	status    = byte
	addr_type = byte
)

const (
	version byte   = 0x05
	cmd     byte   = 0x01 // stream connection
	null    byte   = 0x00 // null byte
	rsv     byte   = 0x00
	granted status = 0x00 // Request granted
)

var (
	// supported auth methods
	NO_AUTH auth = 0x00
	// uname_pwd           auth = 0x02
	auth_not_acceptable auth = 0xFF

	auth_enum   = map[auth]string{}
	status_enum = map[status]string{
		0x00: "request granted",
		0x01: "general failure",
		0x02: "connection not allowed by ruleset",
		0x03: "network unreachable",
		0x04: "host unreachable",
		0x05: "connection refused by destination host",
		0x06: "TTL expired",
		0x07: "command not supported / protocol error",
		0x08: "address type not supported",
	}
)

// type Auth struct {
// 	Username string
// 	Password string
// }

type Socks5Client struct {
	// auth *Auth
	client.Client
}

func New(target client.Context, proxy client.Context) (client_ *Socks5Client, err error) {
	// godoc states:
	// Deferred functions may read and assign to the returning function’s named return values.
	defer func() {
		panicErr := recover()
		if panicErr != nil {
			err = panicErr.(error)
		}
	}()

	if accepted := client.IsIPV4(target.Resolver.(net.IP), proxy.Resolver.(net.IP)); !accepted {
		return nil, err
	}

	client_ = new(Socks5Client)
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
func (c *Socks5Client) Connect(ctx context.Context) error {
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
		go c.tunnel()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c.Worker:
		close(c.Worker)
		return err
	}
}

func (c *Socks5Client) tunnel() {
	var err error
	// shallow copy
	defer func(sh_clone *Socks5Client, err_ *error) {
		if *err_ != nil {
			sh_clone.Close()
		}

		panicErr := recover()
		if panicErr != nil {
			*err_ = errors.New(panicErr.(string))
		}

		sh_clone.Worker <- *err_
	}(c, &err)

	n, err := c.Write([]byte{version, uint8(1), NO_AUTH})
	if err != nil || !(n > 0) {
		if !(n > 0) {
			err = client.ErrWriteTooSmall
		}
		return
	}

	PACKET := make([]byte, 2)
	n, err = c.Read(PACKET)

	if err != nil || !(n > 0) {
		if !(n > 0) {
			err = client.ErrReplyToSmall
		}
		return
	}

	if !(PACKET[0] == version && PACKET[1] == NO_AUTH) {
		if PACKET[1] == auth_not_acceptable {
			err = client.ErrNotAcceptable
		} else {
			err = client.ErrServerChoiceFailure
		}

		return
	}

	PORT := make([]byte, 2)
	{
		binary.BigEndian.PutUint16(PORT, uint16(c.Target.Port))
	}

	PACKET = []byte{}

	PACKET = append(PACKET, version)
	PACKET = append(PACKET, 0x01)
	PACKET = append(PACKET, null)
	PACKET = append(PACKET, 0x01)
	PACKET = append(PACKET, c.Target.Resolver.(net.IP).To4()...)
	PACKET = append(PACKET, PORT...)

	PORT = nil // gc

	n, err = c.Write(PACKET)
	if err != nil || !(n > 0) {
		if !(n > 0) {
			err = client.ErrWriteTooSmall
		}
		return
	}

	PACKET = make([]byte, 2)
	n, err = c.Read(PACKET)

	if err != nil || !(n > 0) {
		if !(n > 0) {
			err = client.ErrReplyToSmall
		}
		return
	}

	if PACKET[1] != granted {
		err = errors.New(status_enum[PACKET[1]])
	}
}
