// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/egor-denisov/wallet-infotecs/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type (
	// Wallet - usecase interfaces.
	Wallet interface {
		CreateNewWalletWithDefaultBalance(c context.Context) (*entity.Wallet, error)
		SendFunds(c context.Context, from string, to string, amount float64) error
		GetWalletHistoryById(c context.Context, walletId string) ([]entity.Transaction, error)
		GetWalletById(c context.Context, walletId string) (*entity.Wallet, error)
	}

	// WalletRepo - repository interfaces.
	WalletRepo interface {
		CreateNewWallet(с context.Context, wallet *entity.Wallet) (*entity.Wallet, error)
		SendFunds(ctx context.Context, transaction *entity.Transaction) error
		GetWalletHistoryById(c context.Context, walletId string) ([]entity.Transaction, error)
		GetWalletById(c context.Context, walletId string) (*entity.Wallet, error)
	}
)
