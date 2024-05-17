package proxifier

import "errors"

var (
	// Err: not an ipv4 address
	ErrNotIPV4 error = errors.New("not ipv4 address")
	// Err: not one of IPV4 or IPV6
	ErrInvalidIPAddr error = errors.New("not one of IPV4 or IPV6")
	// Err: failed writing header packet(s)
	ErrHeaderWrite error = errors.New("failed writing header packet(s)")
	// Err: authentication failed
	ErrAuthFailed error = errors.New("authentication failed")
	// Err: ATYP not supported
	ErrATYP error = errors.New("ATYP not supported")
	// Err: domain name not resolveable
	ErrDomain error = errors.New("domain name not resolveable")
	// Err: not a HTTPS proxy
	ErrNotHTTPSProxy error = errors.New("not a HTTPS proxy")

	// Err: to big can be max of 255 in length
	ErrToBigMax255 error = errors.New("to big can be max of 255 in length")
)
