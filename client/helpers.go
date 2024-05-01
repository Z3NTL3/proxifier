package client

import (
	"context"
	"net"
	"slices"
	"time"
)

// default timeout for LookupHost
var DefaultTimeout_Host time.Duration = time.Second * 5

/*
Looks up the domain, very useful for input check before
providing SOCKS version 5 or 4a with a domain (when you use ATYPE domain)
*/
func LookupHost(input string) (addr []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout_Host)
	defer cancel()

	complete := make(chan int)
	go func() {
		addr, err = net.DefaultResolver.LookupHost(ctx, input)
		complete <- 1
	}()

	select {
		case <-complete:
		case <-ctx.Done():
			err = ctx.Err()
	}
	return
}

/*
Check if given IP address is version 4
*/
func IsIPV4(inputs ...net.IP) bool {
	for _, input := range inputs {
		if input.To4() == nil {
			return false
		}
	}
	return true
}

/*
Check if given IP address is version 6
*/
func IsIPV6(inputs ...net.IP) bool {
	for _, input := range inputs {
		if input.To16() == nil {
			return false
		}
	}
	return true
}

/*
Wrapper to ease IP version 4 and 6 check
*/
func IsIP(inputs ...net.IP) bool {
	check := []bool{IsIPV4(inputs...), IsIPV6(inputs...)}
	return slices.Contains(check, true)
}

/*
Wrapper to ease LookupHost check
*/
func IsDomain(input string) bool {
	addr, err := LookupHost(input)
	if err != nil {
		return false
	}

	if len(addr) > 0 {
		return true
	}

	return false
}

/*
Check whether input is IP version 4,6 or given input is a host
*/
func IsAccepted(inputs ...any) bool {
	for _, input := range inputs {
		ip, isIp := input.(net.IP)

		if !isIp || !IsIP(ip) {
			domain, isDomain := input.(string)

			if !IsDomain(domain) || !isDomain {
				return false
			}
			continue
		}
	}

	return true
}


func MinChar(username, password string) bool {
	return (len(username) > 0 && len(username) <= 255 || 
			len(password) > 0 && len(password) <= 255)
}