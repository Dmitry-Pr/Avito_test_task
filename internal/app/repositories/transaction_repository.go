// Package repositories Description: Этот файл содержит репозиторий для транзакций.
package repositories

import (
	"merch-shop/internal/app/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source=transaction_repository.go -destination=../../../mocks/repositories/transaction_repository.go

const (
	// BuyType описывает тип транзакции покупки.
	BuyType = "buy"
	// TransferType описывает тип транзакции перевода.
	TransferType = "transfer"
)

// TransactionRepositoryInterface описывает репозиторий для транзакций.
type TransactionRepositoryInterface interface {
	Create(tx *gorm.DB, t *models.Transaction) error
	GetTransactionsByUser(tx *gorm.DB, userID uint) ([]models.Transaction, error)
}

// TransactionRepository репозиторий для транзакций.
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository создает новый репозиторий для транзакций.
func NewTransactionRepository(db *gorm.DB) TransactionRepositoryInterface {
	return &TransactionRepository{db: db}
}

// Create создает транзакцию.
func (r *TransactionRepository) Create(tx *gorm.DB, t *models.Transaction) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(t).Error
}

// GetTransactionsByUser получает транзакции по пользователю.
func (r *TransactionRepository) GetTransactionsByUser(tx *gorm.DB, userID uint) ([]models.Transaction, error) {
	if tx == nil {
		tx = r.db
	}
	var transactions []models.Transaction
	err := tx.Where("user_id = ? OR to_user_id = ?", userID, userID).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}
