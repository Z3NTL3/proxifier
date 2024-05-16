package client

import (
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"time"

	uri "net/url"
)

type HTTPClient struct {
	Auth
	TLS bool
}

type Conn  = net.Conn


/*
	Can perform HTTP forward and tunnel requests with to given proxy server.
	TLS is ensured on HTTP TUNNEL requests.

	Reference of implementation:
	- https://en.wikipedia.org/wiki/Proxy_server#Implementations_of_proxies
	- https://en.wikipedia.org/wiki/HTTP_tunnel#HTTP_CONNECT_method
*/
func (c *HTTPClient) PROXY(url string, proxy Context, timeout time.Duration) (conn Conn, err error) {
	defer func(){
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}

		if err != nil {
			conn.Close()
		}
	}()

	// not valid ipv 4 or 6
	if !IsIPV4(proxy.Resolver.(net.IP)) { // may panic
		return nil, ErrNotValidIP
	}
	
	URI, err := uri.Parse(url); if err != nil {
		return
	}

	var auth string

	if len(c.Auth.Username) > 0 || len(c.Auth.Password) > 0 {
		raw_auth_typ := make([]byte, 1)
		{
			raw_auth_typ = append(raw_auth_typ, []byte(c.Auth.Username)...)
			raw_auth_typ = append(raw_auth_typ, byte(':'))
			raw_auth_typ = append(raw_auth_typ, []byte(c.Auth.Password)...)
		}

		auth = base64.StdEncoding.EncodeToString(raw_auth_typ)
	}

	conn, err = net.DialTCP("tcp", nil, &net.TCPAddr{
		IP: proxy.Resolver.(net.IP),
		Port: proxy.Port,
	}); if err != nil {return}


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

	if _, err = conn.Write(PACKET); err != nil {return}

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