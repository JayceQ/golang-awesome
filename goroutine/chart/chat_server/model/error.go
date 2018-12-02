package model

import "errors"

var (
	ErrUserNotExist  = errors.New("user not exist")
	ErrInvalidPasswd = errors.New("passwd or username not right")
	ErrInvalidParams = errors.New("invalid params")
	ErrUserExist     = errors.New("user exist")
)
