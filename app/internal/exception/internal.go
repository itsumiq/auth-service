package exception

import "errors"

var (
	ErrNotFound       = errors.New("Not found entry")
	ErrDuplicateEntry = errors.New("Dublicate entry")
	ErrInternal       = errors.New("Internal error")
	ErrVerify         = errors.New("Verification error")
	ErrTokenExpired   = errors.New("Token expired")
)
