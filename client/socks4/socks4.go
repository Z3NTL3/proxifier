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

	GRANTED           reply = 0x5A // Request granted
	REJECTED_FAILED   reply = 0x5B // Request rejected or failed
	IDENT_UNREACHABLE reply = 0x5C // Request failed because client is not running identd (or not reachable from server)
	IDENT_UNVERIFIED  reply = 0x5D // Request failed because client's identd could not confirm the user ID in the request
)

type Socks4Client struct {
	*net.TCPConn
	target client.TargetCtx
	proxy  client.ProxyCtx
	worker chan struct {
		err error
	}
}

/*
Creates a new SOCKS4 Connect client
*/
func New(target client.TargetCtx, proxy client.ProxyCtx) *Socks4Client {
	return &Socks4Client{
		target: target,
		proxy:  proxy,
		worker: make(chan struct {
			err error
		}, 2),
	}
}

/*
Connects to the target and tunnels through proxy
*/
func (c *Socks4Client) Connect(uid []byte, ctx context.Context) error {
	has_null := false // has termination byte ? (required)
	for _, b := range uid {
		if b == 0x00 {
			has_null = true
		}
	}

	if !has_null {
		uid = append(uid, 0x00)
	}

	go func(cp_client *Socks4Client, cp_uid []byte,
	) {
		proxy_addr := net.TCPAddr{
			IP:   net.ParseIP(cp_client.proxy.IP),
			Port: cp_client.proxy.Port,
		}

		conn, err := net.DialTCP("tcp", nil, &proxy_addr)
		if err != nil {
			cp_client.worker <- struct {
				err error
			}{err: err}
			return
		}

		c.TCPConn = conn

		go c.sent_packet(cp_uid)
	}(c, uid)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case v := <-c.worker:
		close(c.worker)
		return v.err
	}
}

func (c *Socks4Client) sent_packet(uid []byte) {
	var err error
	defer func(client *Socks4Client, err_ *error) {
		panicErr := recover()
		if panicErr != nil {
			err = errors.New(panicErr.(string))
		}

		client.worker <- struct {
			err error
		}{
			err: *err_,
		}
	}(c, &err)

	IP := net.ParseIP(c.target.IP).To4()
	if IP == nil {
		err = errors.New("IPV4 parse failure")
		return
	}

	var HEADER []byte
	{
		PORT := make([]byte, 2)
		binary.BigEndian.PutUint16(PORT, uint16(c.target.Port))

		HEADER = append(HEADER, VERSION)
		HEADER = append(HEADER, CMD)
		HEADER = append(HEADER, PORT...)
		HEADER = append(HEADER, IP.To4()...)
		HEADER = append(HEADER, uid...)
	}

	if n, err := c.Write(HEADER); err != nil || !(n > 0) {
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
	case REJECTED_FAILED:
		err = errors.New("Request rejected or failed")
	case IDENT_UNREACHABLE:
		err = errors.New("Request failed because client is not running identd (or not reachable from server)")
	case IDENT_UNVERIFIED:
		err = errors.New("Request failed because client's identd could not confirm the user ID in the request")

	default:
		err = errors.New("unknown reply byte received from proxy")
	}
}
