package model

import "errors"

// Typed errors for downstream / business failures in the checkout flow.
var (
	ErrValidation        = errors.New("validation failed")
	ErrInventory         = errors.New("inventory failed")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrPayment           = errors.New("payment failed")
)
