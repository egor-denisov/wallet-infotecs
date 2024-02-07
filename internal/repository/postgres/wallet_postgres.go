package repo

import (
	"context"
	"fmt"

	"github.com/egor-denisov/wallet-infotecs/internal/entity"
	"github.com/egor-denisov/wallet-infotecs/pkg/postgres"
)

// WalletRepo -.
type WalletRepo struct {
	*postgres.Postgres
}

// NewWalletRepo -.
func NewWalletRepo(pg *postgres.Postgres) *WalletRepo {
	return &WalletRepo{pg}
}

// CreateNewWallet - creating new wallet entry  in the db.
func (r *WalletRepo) CreateNewWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error) {
	_, err := r.DB.Model(wallet).
		Insert()

	if err != nil {
		return nil, fmt.Errorf("WalletRepo - CreateNewWallet - r.DB: %w", err)
	}
	return wallet, nil
}

// SendFunds - decreasing the balance of the sender and an increasing the receiver. Adding an entry to a transaction table.
func (r *WalletRepo) SendFunds(ctx context.Context, transaction *entity.Transaction) error {
	// Using the db transaction
	tx, err := r.DB.Begin()
    if err != nil {
        return fmt.Errorf("WalletRepo - SendFunds - r.DB: %w", err)
    }

	defer tx.Close()
	// Decreasing the balance of the sender
	res, err := r.DB.Model(&entity.Wallet{}).
		Set("balance = balance - ?", transaction.Amount).
		Where("id = ?", transaction.From).
		Update()
	// If error then rollback the transaction
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("WalletRepo - SendFunds - r.DB: %w", err)
	}
	// If walletId is not found then return 404
	if res.RowsAffected() == 0 {
		tx.Rollback()
		return entity.ErrWalletNotFound
	}
	// Increasing the balance of the receiver
	res, err = r.DB.Model(&entity.Wallet{}).
		Set("balance = balance + ?", transaction.Amount).
		Where("id = ?", transaction.To).
		Update()
	// If error or walletId is not found then rollback the transaction
	if err != nil || res.RowsAffected() == 0 {
		tx.Rollback()
		return fmt.Errorf("WalletRepo - SendFunds - r.DB: %w", err)
	}
	// Adding an entry to a transaction table
	_, err = r.DB.Model(transaction).
		Insert()
	
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("WalletRepo - SendFunds - r.DB: %w", err)
	}
	// Make commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("WalletRepo - SendFunds - r.DB: %w", err)
	}
	return nil
}

// GetWalletHistoryById - getting all transaction records from the user with the walletId.
func (r *WalletRepo) GetWalletHistoryById(ctx context.Context, walletId string) ([]entity.Transaction, error) {
	transactions := make([]entity.Transaction, 0)
	// If walletId is not found or error, return error
	count, err := r.DB.Model(&entity.Wallet{}).
		Where("id = ?", walletId).
		SelectAndCount()
	
	if err != nil || count == 0 {
		return nil, fmt.Errorf("WalletRepo - GetWalletHistoryById - r.DB: %w", err)
	}

	err = r.DB.Model(&transactions).
		Where("from_wallet_id = ?", walletId).
		Select()
		
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetWalletHistoryById - r.DB: %w", err)
	}
	return transactions, nil
}

// GetWalletById - getting wallet info by walletId.
func (r *WalletRepo) GetWalletById(ctx context.Context, walletId string) (*entity.Wallet, error) {
	wallet := new(entity.Wallet)
	err := r.DB.Model(wallet).
		Where("id = ?", walletId).
		Select()
		
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetWalletById - r.DB: %w", err)
	}
	return wallet, nil
}