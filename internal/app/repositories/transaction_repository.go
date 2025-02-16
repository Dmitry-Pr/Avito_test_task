package repositories

import (
	"merch-shop/internal/app/models"

	"gorm.io/gorm"
)

const (
	BuyType      = "buy"
	TransferType = "transfer"
)

type TransactionRepositoryInterface interface {
	Create(tx *gorm.DB, t *models.Transaction) error
	GetTransactionsByUser(tx *gorm.DB, userID uint) ([]models.Transaction, error)
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepositoryInterface {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *gorm.DB, t *models.Transaction) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(t).Error
}

func (r *TransactionRepository) GetTransactionsByUser(tx *gorm.DB, userID uint) ([]models.Transaction, error) {
	if tx == nil {
		tx = r.db
	}
	var transactions []models.Transaction
	err := tx.Where("user_id = ? OR to_user_id = ?", userID, userID).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}
