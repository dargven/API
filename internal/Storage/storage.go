package storage

import "errors"

var (
	ErrURLNotFound   = errors.New("url not found")
	ErrURLExists     = errors.New("url exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user with this email already exists")
	ErrEventNotFound = errors.New("event not found")
)
