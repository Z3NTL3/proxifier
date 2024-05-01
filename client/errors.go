package client

import "errors"

var (
	ErrNotValidIP error = errors.New("not a valid ipv 4 or 6")
	ErrUnsupported error = errors.New("unsupported input")
	ErrNotValidClient error = errors.New("not a valid client")
	ErrHeaderWrite error = errors.New("failed writing header packet(s)")
	ErrAuthFailed error = errors.New("authentication failed")
	ErrMax255Char error = errors.New("can be maximum of 255 charachters")
	ErrATYP error = errors.New("ATYP not supported")
	ErrNotValidDomain error = errors.New("not valid domain")
)
