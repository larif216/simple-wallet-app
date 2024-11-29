package entity

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrDisbursementNotFound = errors.New("disbursement not found")
	ErrInsufficientBalance  = errors.New("insufficient balance")
)
