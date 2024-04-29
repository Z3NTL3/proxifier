package client

import "net"

func IsIPV4(inputs ...net.IP) error {
	for _, input := range inputs {
		if input.To4() == nil {
			return ErrNotIpv4
		}
	}
	return nil
}
