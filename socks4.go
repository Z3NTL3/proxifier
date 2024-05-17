package proxifier

import (
	"encoding/binary"
	"errors"
	"net"
)

type (
	reply = byte // for convenience
)

const (
	SOCKS4 byte = 0x04
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

func (c *Socks4Client) init(target, proxy *Context) (err error) {
	// according to go doc
	// defer func may assign to named returns
	defer func() {
		panicErr := recover()
		if panicErr != nil {
			err = panicErr.(error)
			// may also panic but its always an error type so
			// in our context
		}
	}()

	// not valid IPV4
	if !IsIPV4(target.Resolver.(net.IP), proxy.Resolver.(net.IP)) { // may panic
		return ErrNotIPV4
	}

	c.Client =  Client{
		target: target,
		proxy:  proxy,
		worker: make(chan error),
	}
	c.UID = UID_NULL
	
	return nil
}

func (c *Socks4Client) setup() chan error {
	has_null := false // has termination byte ? (required)
	for _, b := range c.UID {
		if b == NULL {
			has_null = true
		}
	}

	if !has_null {
		c.UID = append(c.UID, NULL)
	}
	
	go func() {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   c.proxy.Resolver.(net.IP),
			Port: c.proxy.Port,
		})
		if err != nil {
			c.worker <- err
			return
		}

		c.TCPConn = conn
		go c.tunnel(c.UID)
	}()

	return c.worker
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

		sh_clone.worker <- *err_
	}(c, &err)

	var PACKET []byte
	{
		PORT := make([]byte, 2)
		binary.BigEndian.PutUint16(PORT, uint16(c.target.Port))

		PACKET = append(PACKET, SOCKS4)
		PACKET = append(PACKET, CMD)
		PACKET = append(PACKET, PORT...)
		PACKET = append(PACKET, c.target.Resolver.(net.IP).To4()...)
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