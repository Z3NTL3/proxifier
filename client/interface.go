package client

import (
	"errors"
	"net"
)

type (
	Context struct {
		IP   net.IP
		Port int
	}

	ProxyCtx  = Context
	TargetCtx = Context
)

const NOT_IPV4 = "not a valid ipv4 address"

func IsIPV4(target net.IP, proxy net.IP) error {
	if target.To4() == nil || proxy.To4() == nil {
		return errors.New(NOT_IPV4)
	}
	return nil
}
