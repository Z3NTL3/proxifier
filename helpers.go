package proxifier

import (
	"context"
	"net"
	"slices"
	"time"
)

// Default timeout for LookupHost.
// Can be changed if desired
var DefaultTimeout_Host time.Duration = time.Second * 5

// Resolves domain. Used for input validation
func LookupHost(input string) (addr []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout_Host)
	defer cancel()

	return net.DefaultResolver.LookupHost(ctx, input)
}

// Reports whether ``inputs`` contain any invalid IPV4 address
func IsIPV4(inputs ...net.IP) bool {
	for _, input := range inputs {
		if input.To4() == nil {
			return false
		}
	}
	return true
}

// Reports whether ``inputs`` contain any invalid IPV6 address.
func IsIPV6(inputs ...net.IP) bool {
	for _, input := range inputs {
		if input.To16() == nil {
			return false
		}
	}
	return true
}

// Reports whether ``inputs`` contain any valid IPV4/6 address.
func IsIP(inputs ...net.IP) bool {
	check := []bool{IsIPV4(inputs...), IsIPV6(inputs...)}
	return slices.Contains(check, true)
}

// Wrapper for LookupHost. With additional validation.
func ValidateDomain(input string) bool {
	if !Max255(input) {return false}

	addr, err := LookupHost(input)
	if err != nil {
		return false
	}

	if len(addr) > 0 {
		return true
	}

	return false
}

// Reports whether one of IPV4/6 or host.
func IsAccepted(inputs ...any) bool {
	for _, input := range inputs {
		ip, isIp := input.(net.IP)

		if !isIp || !IsIP(ip) {
			domain, isDomain := input.(string)

			if !ValidateDomain(domain) || !isDomain {
				return false
			}
			continue
		}
	}

	return true
}

// Reports whether one of inputs exceed max length of 255
func Max255(inputs... string) bool {
	for _, input := range inputs {
		if len(input) > 255 {return false}
	}
	return true
}