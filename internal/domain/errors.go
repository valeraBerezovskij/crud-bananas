package domain

import "errors"

var (
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrUserNotFound = errors.New("user with such credentials not found")
)