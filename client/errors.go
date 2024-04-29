package client

import "errors"

var (
	ErrNotIpv4 error = errors.New("not a valid ipv4 address")
	ErrNetIP   error = errors.New("should provide net.IP for this client")
)
