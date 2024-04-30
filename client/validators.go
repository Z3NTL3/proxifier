package client

import (
	"net"
)

func IsIPV4(inputs ...net.IP) bool {
	for _, input := range inputs {
		if input.To4() == nil {
			return false
		}
	}
	return true
}

/*

can add ipv6/domain support for socks5 later.
that is because the main innovation reason of this client

is only to talk with ipv4 socks4/socks5 proxies

TYPE
	type of the address. One of:
	0x01: IPv4 address
	0x03: Domain name
	0x04: IPv6 address
*/

// func IsIPV6(inputs ...net.IP) bool {
// 	for _, input := range inputs {
// 		if input.To16() == nil {
// 			return false
// 		}
// 	}
// 	return true
// }

// func IsDomain(inputs ...string) bool {
// 	for _, input := range inputs {
// 		if !is_domain.MatchString(input) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func IsIP(inputs ...net.IP) bool {
// 	check := []bool{IsIPV4(inputs...), IsIPV6(inputs...)}
// 	return slices.Contains(check, true)
// }

// func IsAccepted(inputs ...any) bool {
// 	for _, input := range inputs {
// 		ip, isIp := input.(net.IP)

// 		if !isIp || !IsIP(ip) {
// 			domain, isDomain := input.(string)

// 			if !IsDomain(domain) || !isDomain {
// 				return false
// 			}
// 			continue
// 		}
// 	}

// 	return true
// }
