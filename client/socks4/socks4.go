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
	reply = byte
)

const (
	VERSION byte = 0x04
	CMD     byte = 0x01

	NULL byte = 0x00

	GRANTED           reply = 0x5A // Request granted
	REJECTED_FAILED   reply = 0x5B // Request rejected or failed
	IDENT_UNREACHABLE reply = 0x5C // Request failed because client is not running identd (or not reachable from server)
	IDENT_UNVERIFIED  reply = 0x5D // Request failed because client's identd could not confirm the user ID in the request
)

var reply_enum = map[reply]string{
	0x5A: "Request granted",
	0x5B: "Request rejected or failed",
	0x5C: "Request failed because client is not running identd (or not reachable from server)",
	0x5D: "Request failed because client's identd could not confirm the user ID in the request",
}

type Socks4Client struct {
	*net.TCPConn
	target client.Context
	proxy  client.Context
	worker chan error
}

/*
Creates a new SOCKS4 Connect client
*/
func New(target client.Context, proxy client.Context) (client_ *Socks4Client, err error) {
	defer func() {
		panicErr := recover()
		if panicErr != nil {
			err = panicErr.(error)
		}
	}()

	// socks4 client requires net.ip
	resolvers := []net.IP{target.Resolver.(net.IP), proxy.Resolver.(net.IP)}

	if err := client.IsIPV4(resolvers[0], resolvers[1]); err != nil {
		return nil, err
	}

	client_ = new(Socks4Client)
	{
		client_.target = target
		client_.proxy = proxy
		client_.worker = make(chan error)
	}

	return
}

/*
Connects to the target and tunnels through proxy
*/
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
			IP:   c.proxy.Resolver.(net.IP),
			Port: c.proxy.Port,
		})
		if err != nil {
			c.worker <- err
			return
		}

		c.TCPConn = conn
		go c.connection_request(uid)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c.worker:
		close(c.worker)
		return err
	}
}

func (c *Socks4Client) connection_request(uid []byte) {
	var err error
	// shallow copy
	defer func(sh_clone *Socks4Client, err_ *error) {
		panicErr := recover()
		if panicErr != nil {
			*err_ = errors.New(panicErr.(string))
		}

		sh_clone.worker <- *err_
	}(c, &err)

	var HEADER []byte
	{
		PORT := make([]byte, 2)
		binary.BigEndian.PutUint16(PORT, uint16(c.target.Port))

		HEADER = append(HEADER, VERSION)
		HEADER = append(HEADER, CMD)
		HEADER = append(HEADER, PORT...)
		HEADER = append(HEADER, c.target.Resolver.(net.IP).To4()...)
		HEADER = append(HEADER, uid...)

	}

	n, err := c.Write(HEADER)
	if err != nil || !(n > 0) {
		if !(n > 0) {
			err = errors.New("failed sending header packet")
		}
		return
	}

	RESPONSE := make([]byte, 8)
	if _, err = c.Read(RESPONSE); err != nil {
		return
	}

	switch RESPONSE[1] {
	case GRANTED:
		// pass
	default:
		err = errors.New(reply_enum[RESPONSE[1]])
	}
}
