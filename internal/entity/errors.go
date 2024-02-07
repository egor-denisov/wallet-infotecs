package entity

import "errors"

var (
	// Wallet errors
	ErrWalletNotFound = errors.New("wallet not found")
	ErrSenderIsReceiver = errors.New("sender is receiver")
)