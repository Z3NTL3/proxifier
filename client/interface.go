package client

import "net"

type (
	Context struct {
		Resolver interface{}
		Port     int
	}

	Client struct {
		*net.TCPConn
		Target Context
		Proxy  Context
		Worker chan error
	}
)
