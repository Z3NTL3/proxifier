package proxifier

import (
	"encoding/binary"
	"errors"
	"net"
)

type status = byte

const (
	SOCKS5 byte = 0x05
	no_auth byte = 0x00
	uname_pwd byte = 0x02
)

var (
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

func(c *Socks5Client) init(target, proxy *Context) (err error) {
	// according to go doc
	// defer func may assign to named returns
	defer func() {
		panicErr := recover()
		if panicErr != nil {
			err = panicErr.(error)
			// may also panic but i know its always an error type so
		}
	}()

	// not ipv 4,6 or is not a host
	if !IsAccepted(target.Resolver) &&
	 IsIP(proxy.Resolver.(net.IP)){
		err = ErrATYP
	}

	c.Client = Client{
		target: target,
		proxy:  proxy,
		worker: make(chan error),
	}
	c.Auth = Auth{}

	return
}

func (c *Socks5Client) setup() chan error {
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
		go c.tunnel()
	}()
	
	return c.worker
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

		sh_clone.worker <- *err_
	}(c, &err)

	AUTH := no_auth
	if len(c.Username) > 0 || len(c.Password) > 0 {
		AUTH = uname_pwd
	}

	var PACKET []byte = []byte{SOCKS5, uint8(1), AUTH}
	
	n, err := c.Write(PACKET)
	if err != nil || !(n > 0) {
		if !(n > 0) {
			err = ErrHeaderWrite
		}
		return
	}

	PACKET = make([]byte, 2)
	n, err = c.Read(PACKET)
	if err != nil || !(n > 0) {
		if !(n > 0) {
			err = ErrHeaderWrite
		}
		return
	}

	if !(PACKET[0] ==  SOCKS5 && PACKET[1] == AUTH) {
		err = ErrAuthFailed
		return
	}
	
	PACKET = []byte{} // clear

	if AUTH == uname_pwd {
		PACKET = append(PACKET, 0x01)
		PACKET = append(PACKET, uint8(len(c.Username)))
		PACKET = append(PACKET, []byte(c.Username)...)
		PACKET = append(PACKET, uint8(len(c.Password)))
		PACKET = append(PACKET, []byte(c.Password)...)

		n, err = c.Write(PACKET)
		if err != nil || !(n > 0) {
			if !(n > 0) {
				err = ErrHeaderWrite
			}
			return
		}

		PACKET = make([]byte, 2)
		n, err = c.Read(PACKET)
		if err != nil || !(n > 0) {
			if !(n > 0) {
				err = ErrHeaderWrite
			}
			return
		}

		if !(PACKET[0] == 0x01 && PACKET[1] == 0x00) {
			err = ErrAuthFailed
			return
		}
	}
	
	PACKET = make([]byte, 0)
	{
		PACKET = append(PACKET, SOCKS5, CMD, 0x00)
		switch ATYP := c.target.Resolver.(type) {
			case net.IP:
				if IsIPV4(ATYP){
					PACKET = append(PACKET, 0x01) // 0x01: IPv4 address
					PACKET = append(PACKET, ATYP.To4()...)
				} else if IsIPV6(ATYP) {
					PACKET = append(PACKET, 0x04) // 0x04: IPv6 address
					PACKET = append(PACKET, ATYP.To16()...)
				} else {
					err = ErrATYP
					return
				}
			case string:
				PACKET = append(PACKET, 0x03) // 0x03: Domain name
				PACKET = append(PACKET, uint8(len(ATYP)))
				PACKET = append(PACKET, []byte(ATYP)...)
			default:
				err = ErrATYP
				return
		}

		PORT := make([]byte, 2)
		binary.BigEndian.PutUint16(PORT, uint16(c.target.Port))

		PACKET = append(PACKET, PORT...)

		n, err = c.Write(PACKET)
		if err != nil || !(n > 0) {
			if !(n > 0) {
				err = ErrHeaderWrite
			}
			return
		}

		reply := len(PACKET)
		PACKET = make([]byte, reply)

		n, err = c.Read(PACKET)
		if err != nil || !(n > 0) {
			if !(n > 0) {
				err = ErrHeaderWrite
			}
			return
		}

		if !(PACKET[0] == SOCKS5 && PACKET[1] == 0x00) {
			err = errors.New(status_enum[PACKET[1]])
		} 
	}

}