package domain

import "errors"

var (
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
)
