package entity

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrCategoryInvalid     = errors.New("category does not belong to user")
	ErrInvalidCategoryType = errors.New("category type must be income or expense")

	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
	ErrRefreshTokenInvalid = errors.New("refresh token invalid or expired")
)
