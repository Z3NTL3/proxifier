package proxifier_test

import (
	"net"
	"testing"

	"github.com/Z3NTL3/proxifier"
)

// go test -timeout 30s -run ^TestIsAccepted$ github.com/Z3NTL3/proxifier/tests -v
func TestIsAccepted(t *testing.T) {
	inputs := []any{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("FE80:CD00:0000:0CDE:1257:0000:211E:729C"),
		"pix4.dev",
		"koqweioqwir,.dev",
	}

	for _, input := range inputs {
		if !proxifier.IsAccepted(input) {
			t.Logf("%s not accepted", input)
		}
	}
}

// go test -timeout 30s -run ^TestMax255$ github.com/Z3NTL3/proxifier/tests -v
func TestMax255(t *testing.T) {
	var input string
	for range 255 {
		input += "a"
	}

	if proxifier.Max255(input) {
		t.Log("does not exceed max length of 255")
	}

	input += "a"
	if !proxifier.Max255(input) {
		t.Log("does exceed max length of 255")
	}
}

// go test -timeout 30s -run ^TestValidateDomain$ github.com/Z3NTL3/proxifier/tests -v
func TestValidateDomain(t *testing.T) {
	if proxifier.ValidateDomain("pix4.dev") {
		t.Log("pix4.dev valid domain")
	}

	if !proxifier.ValidateDomain("gogo.qwe.go.dev") {
		t.Log("gogo.qwe.go.dev not a valid domain")
	}

}

// go  test -timeout 30s -run ^TestIsIP$ github.com/Z3NTL3/proxifier/tests -v
func TestIsIP(t *testing.T) {
	ips := []net.IP{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("FE80:CD00:0000:0CDE:1257:0000:211E:729C"),
	}

	if proxifier.IsIP(ips...) {
		t.Logf("valid ip addresses")
	}

}

// go test -timeout 30s -run ^TestLookupHost$ github.com/Z3NTL3/proxifier/tests -v
func TestLookupHost(t *testing.T) {
	addr, err := proxifier.LookupHost("pix4.dev")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("found addr: %s", addr[0])
}

// go test -timeout 30s -run ^TestIPV4$ github.com/Z3NTL3/proxifier/tests -v
func TestIPV4(t *testing.T) {
	ip := net.ParseIP("192.168.1.1")
	if proxifier.IsIPV4(ip) {
		t.Logf("%s is ipv4", ip.String())
	} 
}

// go test -timeout 30s -run ^TestIPV6$ github.com/Z3NTL3/proxifier/tests -v
func TestIPV6(t *testing.T) {	
	ip := net.ParseIP("FE80:CD00:0000:0CDE:1257:0000:211E:729C")
	if proxifier.IsIPV6(ip) {
		t.Logf("%s is ipv6", ip.String())
	} 
}