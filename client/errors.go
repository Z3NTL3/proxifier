package client

import "errors"

var (
	ErrNotIpv4 error = errors.New("not a valid ipv4 address")
	// ErrNotIpv6   error = errors.New("not a valid ipv6 address")
	// ErrNotDomain error = errors.New("not a domain name")
	// ErrNotIP     error = errors.New("not a ip address")
	ErrNetIP               error = errors.New("should provide net.IP for this client")
	ErrClientGreeting      error = errors.New("failed sending client greeting packet")
	ErrFailedHeaderPacket  error = errors.New("failed sending header packet")
	ErrNotAcceptable       error = errors.New("no acceptable methods were offered")
	ErrServerChoiceFailure error = errors.New("server choice reply header failure")
	ErrReplyToSmall        error = errors.New("reply is to small")
	ErrAuthFailure         error = errors.New("authentication failed")
	ErrWriteTooSmall       error = errors.New("write is to small")
)
