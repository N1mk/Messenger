package domain

import "errors"

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrUserIsNotInChat = errors.New("user is not in chat")
	ErrNotFound        = errors.New("not found")
	ErrDuplicated      = errors.New("duplicated")
)
