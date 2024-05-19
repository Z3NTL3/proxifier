package proxifier

import (
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"time"

	uri "net/url"
)

// Client for both HTTP and HTTPS proxy
//
// Implements HTTP FORWARD and HTTP TUNNEL mechanism.
type HTTPClient struct {
	Auth // Authentication
	TLS bool // TLS
}

type Conn = net.Conn

// # Doc
//
// Tunnels through given proxy, for HTTPS you should 
// set [proxifier.HTTPClient] to true, which will perform [HTTPS proxy (HTTP TUNNEL)] rather than [HTTP proxy (HTTP FORWARD)]. 
// Defaults to false which is HTTP proxy aka [HTTP proxy (HTTP FORWARD)].
//
// Only HTTP TUNNEL (HTTPS) allows arbitrary data after initial handshake.
//
// [HTTP proxy (HTTP FORWARD)]: https://en.wikipedia.org/wiki/Proxy_server#Implementations_of_proxies
// [HTTPS proxy (HTTP TUNNEL)]: https://en.wikipedia.org/wiki/HTTP_tunnel#HTTP_CONNECT_method
func (c *HTTPClient) PROXY(url string, proxy Context, timeout time.Duration) (conn Conn, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()

	// not valid IPV4/6
	if !IsIPV4(proxy.Resolver.(net.IP)) { // may panic
		return nil, ErrInvalidIPAddr
	}

	URI, err := uri.Parse(url)
	if err != nil {
		return
	}

	var auth string

	if len(c.Auth.Username) > 0 || len(c.Auth.Password) > 0 {
		if !Max255(c.Username, c.Password) {
			err = ErrToBigMax255
			return
		}
		raw_auth_typ := make([]byte, 0)
		{
			raw_auth_typ = append(raw_auth_typ, []byte(c.Auth.Username)...)
			raw_auth_typ = append(raw_auth_typ, byte(':'))
			raw_auth_typ = append(raw_auth_typ, []byte(c.Auth.Password)...)
		}

		auth = base64.StdEncoding.EncodeToString(raw_auth_typ)
	}

	conn, err = net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   proxy.Resolver.(net.IP),
		Port: proxy.Port,
	})
	if err != nil {
		return
	}

	if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return
	}

	PACKET := []byte{}

	if c.TLS {
		PACKET = append(PACKET, []byte(fmt.Sprintf("CONNECT %s:443 HTTP/1.1\r\n", URI.Hostname()))...)
	} else {
		PACKET = append(PACKET, []byte(fmt.Sprintf("GET %s HTTP/1.1\r\n", URI.String()))...)
		PACKET = append(PACKET, []byte(fmt.Sprintf("Host: %s\r\n", URI.Hostname()))...)
		PACKET = append(PACKET, []byte("Connection: close\r\n")...)
	}

	if len(auth) > 0 {
		PACKET = append(PACKET, []byte(fmt.Sprintf("Proxy-Authorization: Basic %s\r\n", auth))...)
	}

	PACKET = append(PACKET, []byte("\r\n")...) // padding

	if _, err = conn.Write(PACKET); err != nil {
		return
	}

	if c.TLS {
		buff := make([]byte, 1042)
		if _, err = conn.Read(buff); err != nil {
			return
		}

		if !strings.Contains(string(buff), "200") {
			err = ErrNotHTTPSProxy
		}
	}

	return
}
