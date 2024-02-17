package entity

import "errors"

var (
	// Wallet errors
	ErrWalletNotFound = errors.New("wallet not found")
	ErrWrongAmount = errors.New("wrong amount")
	ErrSenderIsReceiver = errors.New("sender is receiver")
)