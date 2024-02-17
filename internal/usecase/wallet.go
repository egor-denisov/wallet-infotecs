package usecase

import (
	"context"
	"fmt"

	"github.com/egor-denisov/wallet-infotecs/internal/entity"
)

// WalletUseCase -.
type WalletUseCase struct {
	repo   WalletRepo
	DefaultBalance float64
}

// New -.
func New(r WalletRepo, b float64) *WalletUseCase {
	return &WalletUseCase{
		repo:   r,
		DefaultBalance: b,
	}
}

// CreateNewWallet - creating a new wallet
func (w *WalletUseCase) CreateNewWalletWithDefaultBalance(ctx context.Context) (*entity.Wallet, error) {
	// Create a new instance of the wallet with default balance
	defaultWallet := &entity.Wallet{
		Balance: w.DefaultBalance,
	}

	wallet, err := w.repo.CreateNewWallet(ctx, defaultWallet)
	if err != nil {
		return nil, fmt.Errorf("WalletUseCase - CreateNewWalletWithDefaultBalance - w.repo.CreateNewWallet: %w", err)
	}

	return wallet, nil
}

// SendFunds - sending funds between wallets
func (w *WalletUseCase) SendFunds(ctx context.Context, from string, to string, amount float64) error {
	if amount <= 0 {
		return entity.ErrWrongAmount
	}

	transaction := &entity.Transaction{
		From: from,
		To: to,
		Amount: amount,
	}
	if transaction.From == transaction.To {
		return entity.ErrSenderIsReceiver
	}

	err := w.repo.SendFunds(ctx, transaction)
	if err != nil {
		return fmt.Errorf("WalletUseCase - SendFunds - w.repo.SendFunds: %w", err)
	}

	return nil
}

// GetWalletHistoryById - getting a history of a wallet
func (w *WalletUseCase) GetWalletHistoryById(ctx context.Context, walletId string) ([]entity.Transaction, error) {
	transactions, err := w.repo.GetWalletHistoryById(ctx, walletId)
	if err != nil {
		return nil, fmt.Errorf("WalletUseCase - GetWalletHistoryById - w.repo.GetWalletHistoryById: %w", err)
	}
	
	return transactions, nil
}

// GetWalletById - getting a wallet by id
func (w *WalletUseCase) GetWalletById(ctx context.Context, walletId string) (*entity.Wallet, error) {
	wallet, err := w.repo.GetWalletById(ctx, walletId)
	if err != nil {
		return nil, fmt.Errorf("WalletUseCase - GetWalletById - w.repo.GetWalletById: %w", err)
	}
	
	return wallet, nil
}