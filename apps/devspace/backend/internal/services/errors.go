package services

import "errors"

var ErrInternal = errors.New("internal server error")
var ErrUserExists = errors.New("user with such email already exists")
var ErrUserNotFound = errors.New("user with such email or password is not found")
var ErrUnauthorized = errors.New("invalid credentials")
