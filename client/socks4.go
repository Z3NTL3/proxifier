package client

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
)

type (
	reply = byte // for convenience
)

const (
	version byte = 0x04
	granted reply = 0x5A // Request granted
)


var (
	UID_NULL []byte = []byte{NULL} // for convenience

	reply_enum = map[reply]string{
		0x5A: "Request granted",
		0x5B: "Request rejected or failed",
		0x5C: "Request failed because client is not running identd (or not reachable from server)",
		0x5D: "Request failed because client's identd could not confirm the user ID in the request",
	}
)

func (c *Socks4Client) Connect(uid []byte, ctx context.Context) error {
	has_null := false // has termination byte ? (required)
	for _, b := range uid {
		if b == NULL {
			has_null = true
		}
	}

	if !has_null {
		uid = append(uid, NULL)
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

	var PACKET []byte
	{
		PORT := make([]byte, 2)
		binary.BigEndian.PutUint16(PORT, uint16(c.Target.Port))

		PACKET = append(PACKET, version)
		PACKET = append(PACKET, CMD)
		PACKET = append(PACKET, PORT...)
		PACKET = append(PACKET, c.Target.Resolver.(net.IP).To4()...)
		PACKET = append(PACKET, uid...)

	}

	n, err := c.Write(PACKET)
	if err != nil || !(n > 0) {
		if !(n > 0) {
			err = ErrHeaderWrite
		}
		return
	}

	PACKET = make([]byte, 8)
	if _, err = c.Read(PACKET); err != nil {
		return
	}

	if PACKET[1] != granted {
		err = errors.New(reply_enum[PACKET[1]])
	}
}