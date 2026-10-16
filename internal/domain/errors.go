package domain

import "errors"

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUserIsNotInChat   = errors.New("user is not in chat")
	ErrNotFound          = errors.New("not found")
	ErrDuplicated        = errors.New("duplicated")
)
